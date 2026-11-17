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

Published [preview.12](https://github.com/oneworkspacestudio/my-voice/releases/tag/v0.1.0-preview.12)
from 5967de2, with 16 immutable assets verified. The managed ZIP is 88,363,978
bytes, SHA-256 77df0d5f4d49951da98bb4372ada69b9d3988086d9e8136688b8057099c383d8.
Anonymous metadata and archive downloads passed. Workspace ba64784 promotes
the exact assets and pin, retaining unrelated entries. Actions remain disabled.

Fresh supported synthetic NTFS preview.140 setup selected only Accounts/Voice.
The ordinary unauthenticated installer discovered and applied Voice10→12,
verified all 116 files, and retained manager, Accounts and owner identity.
MyAccounts OIDC authorized the private grant. Phonon and selected resident Pocket
produced real PCM, the proxy denied the direct endpoint, and protected rollback
restored Voice10 with the remaining installation intact. Final release desktop
and 390px mobile settings screenshots were reviewed. All synthetic services stopped.

An isolated WonderPal 0.2.3 relay and browser test consumed real resident Pocket
audio from the released binary; first relay PCM was 28–39ms with two threads,
same worker across replies, browser playback/Stop and revocation passed. Artifacts
are under the local wonderpal-voice-options-20261117 tooling folder. Owner
installation/pairing is pending; Android remains deferred per the user's choice.

## Authorized PC2 installation and Hosting pairing

User then requested installation and hosting. Under that authorization, PC2 used
the normal preview.140 Updates surface for a Voice-only 10→12 update. All nine stages
completed, including checksum verification, 116.5 MB checked recovery and
retained recovery `update-20261117T211611-dd668e17e066`. Ordinary cleanup retired
one eligible older automatic recovery (645.2 MB). Comparing 44 installed app
entries shows only Voice changed; the public health endpoint reports all 46
services running. Authenticated served JS/CSS match the accepted release build.
Remote evidence is the supported updater, authenticated assets and actual app
behavior, not an independent remote-filesystem binary audit.

Pocket 18M and MOSS were already installed. The owner approved Qwen 0.6B's licence
and 1 GB download; the normal My Voice installer completed its audio readiness
check and reports Installed/Ready. Existing Qwen17 and saved private voices were
left in place. Pocket 18M is selected, enabled and loaded for WonderPal. A dedicated
grant was saved only as encrypted production Hosting values alongside its
same-PC loopback origin. No credential is in source, captures or these notes.

WonderPal 0.2.3 deployment 96abb38c1316dbf2a0c5ee8265bcb741 rebuilt the retained
verified ZIP with current protected settings and published at 21:27:19 UTC.
Previous production 5cf3a2d77b730d4c073302c52ddd4965 remains recoverable. Live
status advertises Workspace Voice with DeepSeek/read-aloud retained. An actual
live reply and replay returned HTTP 200 mono PCM; first data 152/194 ms after their
speech requests, with complete nonempty streams and working Stop. Owner voice
references were not used. The owner Studio tab was not navigated. Temporary
administration tabs/helper were closed; the new live WonderPal page is delivered.
Local privacy-filtered receipts/screenshots are under the tooling task's
`owner-pc2` folder. PC2 delivery is complete; Android and other PCs were not
updated in this request.
