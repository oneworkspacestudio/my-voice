"""Stock English resident engines. No request text or audio is written to disk."""
import ctypes
import http.client
import json
import os
from pathlib import Path
import socket
import struct
import subprocess
import time
import urllib.parse

import numpy as np


class Phonon:
    def __init__(self, model, _threads):
        import ptts
        self.engine = ptts.TTS(str(model), lang="en", quant="q8", voice_catalog="")
        if self.engine.sample_rate != 24000:
            raise RuntimeError("Unexpected sample rate")

    def speak(self, text, emit):
        with self.engine.stream(text) as audio:
            for samples in audio:
                emit(samples, 24000)


class Pocket:
    def __init__(self, model, _threads):
        runtime = Path(__file__).resolve().parent.parent
        # The native server is an internal child in the apphost Job Object. It
        # receives no grant/token and is bound to a random loopback-only port.
        for executable in ("pocket-tts.exe", "pocket-tts-portable.exe"):
            with socket.socket() as reservation:
                reservation.bind(("127.0.0.1", 0))
                port = reservation.getsockname()[1]
            self.process = subprocess.Popen([
                str(runtime/executable), "serve", "--host", "127.0.0.1", "--port", str(port),
                "--config", str(model/"config.yaml"), "--default-voice", str(model/"alba.safetensors"),
            ], stdin=subprocess.DEVNULL, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
                creationflags=subprocess.CREATE_NO_WINDOW,
                env={**os.environ, "FIRST_CHUNK_LENGTH_SECONDS": "0"})
            self.port = port
            for _ in range(1200):
                if self.process.poll() is not None:
                    break
                conn = http.client.HTTPConnection("127.0.0.1", port, timeout=0.2)
                try:
                    conn.request("GET", "/health")
                    response = conn.getresponse()
                    if response.status == 200:
                        response.read()
                        return
                except OSError:
                    pass
                finally:
                    conn.close()
                time.sleep(0.05)
            if self.process.poll() is None:
                self.process.kill()
                self.process.wait()
                raise RuntimeError("Pocket did not start")
            if (self.process.returncode & 0xffffffff) != 0xc000001d:
                raise RuntimeError("Pocket failed to start")
        raise RuntimeError("No compatible Pocket worker")

    def speak(self, text, emit):
        conn = http.client.HTTPConnection("127.0.0.1", self.port, timeout=40)
        try:
            conn.request("POST", "/tts", urllib.parse.urlencode({"text": text}),
                         {"Content-Type": "application/x-www-form-urlencoded"})
            response = conn.getresponse()
            if response.status != 200:
                raise RuntimeError("Pocket could not speak")
            header = response.read(44)
            if (len(header) != 44 or header[:4] != b"RIFF" or header[8:16] != b"WAVEfmt "
                    or header[36:40] != b"data" or struct.unpack_from("<HH", header, 20) != (1, 1)
                    or struct.unpack_from("<H", header, 34)[0] != 16):
                raise RuntimeError("Unexpected Pocket WAV")
            rate = struct.unpack_from("<I", header, 24)[0]
            pending = b""
            while data := response.read1(8192):
                data = pending + data
                count = len(data) & ~1
                emit(np.frombuffer(data[:count], dtype="<i2").astype(np.float32)/32768, rate)
                pending = data[count:]
            if pending:
                raise RuntimeError("Truncated Pocket PCM")
        finally:
            conn.close()


class Moss:
    def __init__(self, model, threads):
        from ort_cpu_runtime import OrtCpuRuntime
        import sentencepiece as spm
        self.runtime = OrtCpuRuntime(str(model), thread_count=threads, max_new_frames=768,
                                    execution_provider="cpu", sample_mode="fixed")
        self.tokenizer = spm.SentencePieceProcessor(model_file=str(
            self.runtime.resolve_manifest_relative_path(self.runtime.manifest["model_files"]["tokenizer_model"])))
        state = json.loads((model/"preset-state.json").read_text(encoding="utf-8"))
        if state["engine"] != "moss" or state["version"] != "f52645cb":
            raise RuntimeError("Incompatible preset")
        self.codes = state["codes"]
        self.rate = int(self.runtime.codec_meta["codec_config"]["sample_rate"])

    def speak(self, text, emit):
        runtime = self.runtime
        runtime.codec_streaming_session.reset()
        rows = runtime.build_voice_clone_request_rows(self.codes, self.tokenizer.encode(text, out_type=int))
        pending = []

        def decode():
            result = runtime.codec_streaming_session.run_frames(pending)
            pending.clear()
            if result is not None:
                audio, length = result
                if length > 0:
                    emit(audio[0, :, :length].mean(axis=0), self.rate)

        def frame(_frames, _step, codes):
            pending.append(codes)
            if len(pending) >= 4:
                decode()

        try:
            frames = runtime.generate_audio_frames(rows, on_frame=frame)
            if len(frames) >= 768:
                raise RuntimeError("Speech limit exceeded")
            if pending:
                decode()
        finally:
            runtime.codec_streaming_session.reset()


class Qwen:
    def __init__(self, model, threads):
        runtime = Path(__file__).resolve().parent.parent/"qwen"
        self.dll_directory = os.add_dll_directory(str(runtime))
        self.library = ctypes.CDLL(str(runtime/"myvoice-qwen.dll"))
        self.library.voice_create.argtypes = [ctypes.c_wchar_p, ctypes.c_char_p, ctypes.c_char_p,
                                             ctypes.c_char_p, ctypes.c_int]
        self.library.voice_create.restype = ctypes.c_void_p
        self.callback_type = ctypes.CFUNCTYPE(ctypes.c_int, ctypes.POINTER(ctypes.c_float),
                                              ctypes.c_size_t, ctypes.c_int32, ctypes.c_void_p)
        self.library.voice_stream.argtypes = [ctypes.c_void_p, ctypes.c_char_p, self.callback_type, ctypes.c_void_p]
        self.library.voice_stream.restype = ctypes.c_int
        self.handle = self.library.voice_create(str(runtime), str(model/"model.gguf").encode("utf-8"),
            str(model/"codec.gguf").encode("utf-8"), str(model/"presets/open-voice.wav").encode("utf-8"), threads)
        if not self.handle:
            self.library.voice_last_error.restype = ctypes.c_char_p
            raise RuntimeError("Qwen could not load: " + self.library.voice_last_error().decode("ascii"))

    def speak(self, text, emit):
        failure = []

        def output(data, count, rate, _user):
            try:
                emit(np.ctypeslib.as_array(data, shape=(count,)), rate)
                return 1
            except Exception as error:
                failure.append(error)
                return 0

        callback = self.callback_type(output)
        result = self.library.voice_stream(self.handle, text.encode("utf-8"), callback, None)
        if failure:
            raise failure[0]
        if result:
            raise RuntimeError("Qwen generation interrupted")


def load_engine(name, model, threads):
    return {"phonon": Phonon, "pocket18": Pocket, "pocket7": Pocket,
            "moss": Moss, "qwen06": Qwen}[name](model, threads)
