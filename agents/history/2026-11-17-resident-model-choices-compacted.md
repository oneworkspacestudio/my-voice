# Workspace-first WonderPal speech choices

User chose Workspace voice service first for Pocket TTS, MOSS and Qwen 0.6B;
Android remains unchanged. My Voice preview.12 adds an owner-controlled selector
under WonderPal read-aloud for existing installed Pocket 18M/7M, MOSS Nano,
Qwen3 TTS 0.6B and legacy Phonon. No WonderPal relay contract change is needed.

Selection is private global preference `wonderpal-speech-model`. Switching
cancels the resident worker and active stream while preserving its scoped
credential. The enabled selection restores at restart; uninstall/reinstall of
the active model is blocked. Larger models require an explicit Models download;
Phonon retains its existing automatic pinned data download. Models are not in
release payloads.

One supervised Python facade retains the selected engine. Pocket uses the pinned
native Rust server as a same-Job-Object child, bound to random loopback only;
that upstream child has no bearer authentication and is not published/proxied.
MOSS retains ONNX sessions and incrementally decodes frames. Qwen uses the small
application C++ adapter against pinned b12084 public C APIs, with eight-frame
flush and persistent model/codec/context. Build via prepare-runtime.py requires
MSVC; it fetches and verifies the existing source pin and compiles the adapter.
No executable downloads occur during startup. Runtime extraction remains hashed.

Native synthetic NTFS tests rendered real repeated audio for Pocket18, MOSS,
Qwen06 and Phonon, same worker reuse, switching, grant continuity, restart and
revocation. Owner/CSRF tests cover invalid/missing choice, model removal and
reinstall protections and old-worker cancellation. Four-thread warm first audio:
Pocket39ms, MOSS299ms, Qwen2.826s, Phonon19ms on i9-11900K. Larger-model speed is
not claimed instant. Browser selector, production build, 19 UI regressions and
Go test/vet passed. Desktop/mobile presentation reviewed; no physical Android
test. A temporary frontend text-encoding regression was corrected before source
publication. Keep UTF-8 explicit in Windows Python file operations.

Release publication, public channel promotion and managed upgrade/rollback
verification remain to be recorded below. Owner installation/pairing is pending;
do not use an owner installation as the test fixture.
