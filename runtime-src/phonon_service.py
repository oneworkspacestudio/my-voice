"""Private resident worker supervised by My Voice. Input and audio are ephemeral."""

from hmac import compare_digest
from http.server import BaseHTTPRequestHandler, HTTPServer
import json
import os
from pathlib import Path
import sys

import numpy as np
from resident_engines import load_engine


class SpeechHandler(BaseHTTPRequestHandler):
    protocol_version = "HTTP/1.1"

    def log_message(self, *_args):
        pass

    def setup(self):
        super().setup()
        self.connection.settimeout(10)

    def do_POST(self):
        self.close_connection = True
        authorization = self.headers.get("Authorization", "").encode("latin-1")
        if self.path != "/speech" or not compare_digest(authorization, secret):
            self.send_error(403)
            return
        try:
            size = int(self.headers.get("Content-Length", "0"))
            if not 0 < size <= 8192:
                raise ValueError()
            text = json.loads(self.rfile.read(size))["input"]
            if not isinstance(text, str) or not 0 < len(text.encode("utf-8")) <= 1600:
                raise ValueError()
        except (ValueError, KeyError, TypeError):
            self.send_error(400)
            return
        try:
            started = False
            sample_count = 0

            def emit(samples, rate):
                nonlocal started, sample_count
                # Engines expose mono float32 at 24 kHz or MOSS 48 kHz. Keep
                # decimation aligned across chunks; never change speech speed.
                if rate not in (24000, 48000):
                    raise RuntimeError("Unexpected sample rate")
                samples = np.asarray(samples).reshape(-1)
                if rate == 48000:
                    source_count = len(samples)
                    samples = samples[sample_count % 2::2]
                    sample_count += source_count
                pcm = (np.clip(samples, -1.0, 32767 / 32768) * 32768).astype("<i2")
                if not pcm.nbytes:
                    return
                if not started:
                    start_response()
                    started = True
                self.wfile.write(f"{pcm.nbytes:x}\r\n".encode("ascii") + pcm.tobytes() + b"\r\n")
                self.wfile.flush()

            def start_response():
                self.send_response(200)
                self.send_header("Content-Type", "audio/pcm")
                self.send_header("Transfer-Encoding", "chunked")
                self.send_header("Connection", "close")
                self.end_headers()
            engine.speak(text, emit)
            if not started:
                self.send_error(502)
                return
            self.wfile.write(b"0\r\n\r\n")
        except (BrokenPipeError, ConnectionResetError, TimeoutError):
            pass
        except Exception:
            # Drop malformed or interrupted output without logging request text.
            try:
                self.connection.shutdown(2)
            except OSError:
                pass


if __name__ == "__main__":
    name = sys.argv[1]
    model, marker = map(Path, sys.argv[2:4])
    threads = int(sys.argv[4])
    secret = ("Bearer " + os.environ.pop("MYVOICE_WORKER_TOKEN")).encode("ascii")
    engine = load_engine(name, model, threads)
    with HTTPServer(("127.0.0.1", 0), SpeechHandler) as server:
        pending = marker.with_suffix(".tmp")
        pending.write_text(str(server.server_port), encoding="ascii")
        pending.replace(marker)
        server.serve_forever()
