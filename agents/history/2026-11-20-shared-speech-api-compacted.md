# Shared speech API

My Voice preview.14 replaces app-specific Settings panels with a generic API
key list. Keys are shown once and stored as hashes. Existing pairs are listed
and individually revocable. The new loopback-only API exposes a read-only model
catalog and stock PCM synthesis with the caller's model, preset and optional
speed/pitch. Generic keys cannot change models or access private studio data.
Older direct endpoints keep deployed callers working. Private owner routes
for issuing app-specific keys were removed.

WonderPal persists its own model on its server. My Podcasts owns its narration
model and host preset/style mapping. Existing deployed credentials can be
used by upgraded callers without rotation. Phonon starts on demand, while
other generic speech requests use the installed stock CPU engines.

Use synthetic NTFS fixtures and the existing 116-payload managed update group.
Runtime payload SHA256 matches preview.13. No owner grants or installations
are part of verification. See docs/API.md.

Preview.14 published from be8ec7e and promoted in Workspace 084ae8d. Native
Windows checks generated real Alba/Marius PCM with two API keys and confirmed
forwarded denial and key revocation isolation. Go tests/vet and 19 frontend
tests passed. Anonymous downloads verified all 116 managed files. Owner
installation and a fresh managed upgrade were not exercised.
