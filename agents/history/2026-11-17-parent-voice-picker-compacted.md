# WonderPal parent-facing resident voice selector

User could not choose TTS in WonderPal's parent settings. Previous delivery
offered selection only inside My Voice; WonderPal's Workspace branch displayed
a connection note. This was a missing product control, not missing PC2 models.

Preview.13 adds direct GET/PATCH /v1/audio/models under the existing same-PC
WonderPal grant. It lists only stock choices and installed/readiness metadata,
changes the saved selection, cancels the old worker and preserves the credential.
Browser/cookie/forwarded/query requests are rejected. No download, private voice,
owner-management or public proxy capability is added. Existing grants still work.

My Voice's selector uses separate draft state, follows external changes when
not edited locally and no longer changes editor fields on status polling.
WonderPal 0.2.4 gates its GET/PUT relay by the existing parent token, keeps the
grant on its server and shows installed choices, readiness and refresh. Selection
applies immediately; older services show an upgrade instruction. Android keeps
its device voices.

Go tests/vet, UI builds, Node authorization/shape/compatibility tests and desktop/
mobile picker cases passed. Native synthetic NTFS acceptance exercised real
Pocket/Phonon switching through the parent UI, persisted selection, approved
audio and Stop; same-worker Pocket first PCM 45/34 ms and revocation passed.
No owner data/voice was used. Artifacts: local tooling task
wonderpal-parent-tts-20261117. Publication, managed update and actual PC2
deployment/parent selection remain to be recorded after verification.
