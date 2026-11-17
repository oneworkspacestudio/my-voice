# Resident Phonon read-aloud for WonderPal

My Voice preview.11 adds an owner-controlled WonderPal credential in Settings. The grant is stored as a hash, separate from the existing Podcasts credential. It starts a supervised Windows process with pinned Gradium Phonon 31M Q8. The existing private model storage owns the verified 65 MB weights. Revocation cancels the current stream and stops the resident worker. An enabled grant restores readiness on service restart. Existing studio workflows, voice references and Podcasts remain separate.

The private direct listener accepts only authenticated loopback `POST /v1/audio/speech`, with stock English and bounded mono 24 kHz PCM. It rejects browser/proxy headers. Python receives a per-process credential and uses a separate ephemeral loopback listener. The worker stays loaded while enabled; raw text and audio aren't persisted.

Isolated Windows NTFS fixture verified real audio repeatedly from the same Python worker, forwarding denial, browser playback through WonderPal and revocation. Warm end-to-end first PCM through the local WonderPal relay was 31–35 ms across three fixture requests, excluding model generation by the chat provider. Go tests and vet passed. Model weights remain outside public releases. Workspace update uses the established Voice package namespace and minimum manager, so no manager or companion application release is needed for the independent update.

The live WonderPal project must use the local origin and a protected owner-issued token on the same PC. No owner installation or grant was used in testing.


Published and verified My Voice preview.11 from 3d8f928. The release has 16 immutable assets, including corresponding sources, no model weights, and the 88,343,015-byte managed ZIP SHA-256 0fc995d77ffa46724969263344ee2f7d2978c7e6a7c733b89be5b67d9f658c08. Workspace main 79c625a pins the package and promoted the exact verified public channel.

An anonymous Workspace preview.140 installer selected Voice preview.10 to preview.11, verified all 116 owned files, and completed the supported upgrade. Synthetic OIDC owner pairing warmed Phonon and produced direct PCM. The Workspace proxy denied the direct token endpoint. Protected rollback restored Voice preview.10 while preserving owner identity, Accounts, manager and available recovery. All synthetic services stopped. PC2 WonderPal 0.2.2 is live, with owner installation and pairing pending.
