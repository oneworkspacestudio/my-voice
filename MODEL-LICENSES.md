# Model and preset licences

| Component | Exact upstream source | Conditions |
|---|---|---|
| Pocket runtime | kyutai-labs/pocket-tts v4.0.1 Windows artifact (CLI3.3.0) | MIT core; bundled third-party components include MPL-2.0; notices and corresponding sources accompany the release |
| Pocket 18.2M / 6.9M preset checkpoints | kyutai/pocket-tts-without-voice-cloning; revisions and SHA256 in models.json | CC BY4.0; attribution to Kyutai; preset-only checkpoint |
| Alba voice | Kyutai tts-voices / Alba MacKenna casual recording and model-specific preset state | CC BY4.0; credit Alba MacKenna and Kyutai |
| MOSS TTS Nano ONNX weights | OpenMOSS-Team/MOSS-TTS-Nano-100M-ONNX f52645cb | Apache-2.0 model card; separately downloaded |
| MOSS audio tokenizer | OpenMOSS-Team/MOSS-Audio-Tokenizer-Nano-ONNX; exact revision/digests in models.json | Apache-2.0 model card; separately downloaded |
| MOSS runtime source | OpenMOSS/MOSS-TTS-Nano 8b7bcc9341b3b4ef3a3a58ba1338a7d85ff133eb | Apache-2.0; upstream ort_cpu_runtime.py is unmodified |
| Qwen3-TTS 0.6B Base Q8 | Qwen/Qwen3-TTS-12Hz-0.6B-Base; mradermacher GGUF conversion 0762b31e | Apache-2.0; pinned talker and codec, separately downloaded |
| Qwen3-TTS 1.7B Base Q8 | Qwen/Qwen3-TTS-12Hz-1.7B-Base; official ggml-org GGUF ca27d74b | Apache-2.0; pinned talker and codec, separately downloaded |
| Qwen native runtime | ggml-org/llama.cpp b12084 / f9448cdd | MIT; LLVM OpenMP and Microsoft app-local release CRT notices in THIRD-PARTY-NOTICES.md |
| Open Voice preset reference | Kyutai tts-voices / voice-donations/0a67.wav; pinned revision/digest in models.json | CC0 voluntary Unmute Voice Donation Project reference; state generated locally for MOSS; private reference input for Qwen |

No model weights or personal references are included in the application release.
Owner-initiated installation acknowledges each model's linked licence. The gated
Pocket cloning checkpoint is **not used, mirrored or redistributed**. A browser
access grant does not imply CLI credentials; this preview needs neither for MOSS nor Qwen.

The MOSS source advertises 20 languages, but its published table explicitly names
19 codes. The app exposes those 19 codes and does not invent the missing entry or
Hindi. Only English is covered by native acceptance here. Upstream named demo
states are not offered. The default uses an explicitly CC0 donated reference.

Native cloning QA uses a CC0 reference from Kyutai's Unmute Voice Donation Project,
not an owner's recording. No claim of an identity match to a human is inferred
from signal validation alone. Quality is subjective; compare the provided source
and generated test sample before relying on a voice.

Sources: https://github.com/kyutai-labs/pocket-tts ·
https://huggingface.co/kyutai/tts-voices ·
https://github.com/OpenMOSS/MOSS-TTS-Nano ·
https://huggingface.co/OpenMOSS-Team/MOSS-TTS-Nano-100M-ONNX ·
https://huggingface.co/OpenMOSS-Team/MOSS-Audio-Tokenizer-Nano-ONNX

Generated speech is synthetic. Disclose that when sharing it, and do not present
it as an original recording or evidence that a person said those words.

Qwen exposes ten upstream languages: Chinese, English, German, Italian, Portuguese,
Spanish, Japanese, Korean, French and Russian. The 0.6B quantization is a pinned
community conversion; the 1.7B quantization is published by ggml-org. Neither
requires a Hugging Face login. Model-specific saved references and state are kept
private. A model comparison creates a separate voice and preserves the original.

Qwen sources: https://github.com/QwenLM/Qwen3-TTS;
https://huggingface.co/mradermacher/Qwen3-TTS-12Hz-0.6B-Base-GGUF;
https://huggingface.co/ggml-org/Qwen3-TTS-12Hz-1.7B-Base-GGUF.
