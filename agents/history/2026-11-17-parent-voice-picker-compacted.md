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

Published preview.13 from de34fe1 with 16 immutable assets, verified anonymous
metadata/archive and public channel promotion 29b9199. Anonymous preview.140
installer selected exact digests; synthetic NTFS upgrade, all 116 files, owner
OIDC, real Phonon/Pocket speech, proxy denial and protected rollback passed.
Its first attempt stopped safely on the fixture's disk-space guard; an archived
and checksum-verified disposable native fixture freed enough room for the retry.

PC2 normal Voice-only 12→13 update completed all nine stages in 4m58s.
Recovery update-20261117T225737-59e95284ce0c retained; 46 services running.
Owner Settings confirms preview.13, existing active grant and resident Pocket18.
MOSS/Qwen06 remain installed. WonderPal 0.2.4 published on existing Hosting as
ba351e994f5ca5d0b0b5ed4936068eaf at 23:05:22 UTC. Public hashes/health and
pairing passed; actual parent selection awaits its user-controlled lock.

Follow-up live parent verification complete: production WonderPal directly
selected MOSS, Qwen06 and Pocket18 through the existing grant. Refresh retains
selection. Qwen06 and restored Pocket18 reported loaded/ready. MOSS selection
saved but was still loading when changed; no new MOSS readiness claim.
Pocket18 remains selected and warm. No credentials were rotated or recorded.
