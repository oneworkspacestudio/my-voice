# Architecture

Pocket's official Windows worker targets x86-64-v3. A bundled native baseline
x86-64 build of the identical pinned source provides a fallback only when Windows
returns STATUS_ILLEGAL_INSTRUCTION (0xC000001D). Both run in the same restricted
Job Object with explicit private model/config paths and scrubbed environment.
Deterministic missing/blocked/incompatible worker failures stop without repeated
launches. Job errors expose only fixed recovery guidance, worker name and numeric
OS status, never command lines, stdout/stderr or speech text. Model/config and
database formats are unchanged by this runtime repair.

`MyVoice.exe` is a native Go server with embedded production React/AUI assets.
It uses the public Workspace `pkg/apphost` for MyAccounts OIDC/PKCE, session and
CSRF checks, owner controls, SQLite migrations, protected NTFS storage and backups.
Managed installations disable local accounts. Every sensitive API verifies the
app session; the normal Workspace proxy also requires its own session.

`SpeechEngine` isolates Pocket's native Rust CLI and MOSS's ONNX CPU worker.
The MOSS worker directly uses the pinned, unmodified upstream `OrtCpuRuntime`;
reference resampling/PCM I/O use NumPy, without PyTorch or torchaudio. Saved MOSS
states contain audio-tokenizer codes tied to the exact engine version.
No embedding is advertised as interchangeable between engines.

Runtimes are embedded at build time and unpacked to a content-addressed private
directory. Each extracted file is hashed against its embedded bytes on first use
of each process. Worker environments omit Workspace/OIDC credentials. Workers
start suspended inside kill-on-close Windows Job Objects. Cancellation and parent
exit terminate descendants. The app supervises one inference pipeline at a time.

Model installation requires an owner gesture and licence confirmation, uses only
committed URLs/sizes/digests, stages `.part` files and creates its installation
marker only after verification. Models are separate from public application
downloads. Normal generation uses local files and an offline Hugging Face setting.

SQLite holds owner-scoped projects, jobs, voice metadata, consent timestamps,
explicit grants, dictionaries and audit events. Random identifiers replace user
filenames. A maximum of 20 jobs is queued globally. Script limits are 200 KB,
1,000 segments and 100 dictionary rules. Audio storage is capped at 10 GB and
an assembled export at 256 MB PCM. Uploads are decoded/validated as audio, with
25 MB and 3–30 second reference limits, RMS and clipping checks.

Jobs persist each finished segment. Interrupted running jobs return to the queue
on restart. Targeted regeneration copies unchanged segment bytes into a new job;
exports are assembled separately. Output WAV files are written by rename from
staging files. Duplicate request IDs are scoped to the user and return the same
job, allowing a lost generation response to be retried without duplicate work.

Pronunciation is exact text replacement, account rules followed by project rules.
Speed/pitch/volume are explicitly post-processing controls. No unsupported SSML,
phoneme control, network audio streaming or progress-to-first-audio is implied.

My Drive uses its actual signed-in browser API and verifies OIDC subject equality
before writing. Exports use a deterministic job-specific path, never overwrite,
and verify the downloaded destination SHA256. Failed transfers remain visibly
unsaved; retry reconciles an existing identical file. Private reference recordings
and voice states are never placed in a broad Drive folder.

The suite registry owns `/voice`, base-port+38, the app's OIDC callback, its
private working directory, every `voice-*` package, selection, recovery, shared
picker and update group. Installer integration7 and `native-voice-v1` prevent
older managers accepting an incomplete native integration.
