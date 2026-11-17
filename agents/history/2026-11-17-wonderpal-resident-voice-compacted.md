# Resident Phonon read-aloud for WonderPal

My Voice preview.11 adds an owner-controlled WonderPal credential in Settings. The grant is stored as a hash, separate from the existing Podcasts credential. It starts a supervised Windows process with pinned Gradium Phonon 31M Q8. The existing private model storage owns the verified 65 MB weights. Revocation cancels the current stream and stops the resident worker. An enabled grant restores readiness on service restart. Existing studio workflows, voice references and Podcasts remain separate.

The private direct listener accepts only authenticated loopback `POST /v1/audio/speech`, with stock English and bounded mono 24 kHz PCM. It rejects browser/proxy headers. Python receives a per-process credential and uses a separate ephemeral loopback listener. The worker stays loaded while enabled; raw text and audio aren't persisted.

Isolated Windows NTFS fixture verified real audio repeatedly from the same Python worker, forwarding denial, browser playback through WonderPal and revocation. Warm end-to-end first PCM through the local WonderPal relay was 31–35 ms across three fixture requests, excluding model generation by the chat provider. Go tests and vet passed. Model weights remain outside public releases. Workspace update uses the established Voice package namespace and minimum manager, so no manager or companion application release is needed for the independent update.

The live WonderPal project must use the local origin and a protected owner-issued token on the same PC. No owner installation or grant was used in testing.
