# Native Windows test report

Date: 3 November 2026. Preview0.1.0-preview.2, Windows x64.
All execution used the current physical Windows PC. No VM, WSL, container,
GPU runtime, cloud speech service or hosted runner was used.

## Host and method

Windows11 Home10.0.22631; Intel Core i9-11900K,8 physical/16 logical cores;
31.9GiB usable RAM. Explicit CPU execution; two configured inference threads.
Programs ran outside the source repository, with synthetic isolated NTFS private
data and accounts. No owner's voice, identity, files or installation was a feature
fixture. Consent references came from Kyutai's explicitly CC0 voice donations.

The launcher embeds its own Python/ONNX and native Pocket/audio runtimes. No
separate user Python installation is required. Weights download separately.
Go tests/vet, the production AUI/TypeScript/Vite build, native API acceptance and
native Chrome UI checks passed. Desktop1600x1000 and emulated390x844 browser
viewports passed; this does not establish a physical-phone/microphone test.

## Measured CPU performance

Wall time includes child startup, segmented inference, postprocessing, joining
and both exports. Peak tree working set sums the app and its child processes;
shared mappings can be counted more than once. Sampling interval was80ms.
These are host measurements, not minimum specifications or low-end projections.

| Model | Case | Audio seconds | Wall seconds | RTF | Peak tree MiB |
|---|---|---:|---:|---:|---:|
| pocket18 | short | 18.987 | 4.621 | 0.243 | 383.5 |
| pocket18 | warm | 4.625 | 1.312 | 0.284 | 360.4 |
| pocket18 | long | 125.529 | 22.174 | 0.177 | 426.1 |
| pocket7 | short | 19.707 | 3.298 | 0.167 | 371.8 |
| pocket7 | warm | 4.705 | 1.29 | 0.274 | 354.0 |
| pocket7 | long | 119.129 | 20.002 | 0.168 | 391.7 |
| moss | short | 28.444 | 39.17 | 1.377 | 1673.6 |
| moss | warm | 13.328 | 22.689 | 1.702 | 1653.0 |
| moss | long | 146.809 | 292.786 | 1.994 | 1640.9 |

Pocket workloads were329 characters (about20 seconds),71 characters and1963
characters (about2 minutes). MOSS workloads were230,71 and1433 characters. Actual audio duration varies by engine/reference; no fixed20-second output
is asserted. RTF below1 means faster than real time. Workers load per segment;
“warm” means a subsequent request with installed/cache-ready assets, not an
always-resident model. Playback starts after the assembled job completes;
network time to first audio is not exposed or measured. Cold weight download
time was not benchmarked. The first Pocket case includes fresh native runtime
preparation. MOSS support prepares only on explicit use; a separate first-use
check produced4.16 seconds in19.97 seconds including its runtime preparation. Existing-data process start to health measured about
1.6seconds; initial runtime extraction is excluded from that health timing.

Performance measurements used the identical pinned inference/runtime payload of
the release candidate. Subsequent request-boundary and migration fixes were
verified with the final executable. Exact measurement provenance is retained in
docs/evidence/windows-native.json.

## Native functional acceptance

- Real Pocket18M/7M preset synthesis and MOSS permissioned-reference cloning;
  new text generated from saved model-specific state; WAV/MP3 decoding/duration
  and non-silence checked. Licensed CC0 Open Voice default prepared locally.
- Dictionaries, per-project overrides, distinct-reference two-voice dialogue,
  pauses, long-text ordering/chunking and bulk ZIP passed.
- Selected-segment regeneration preserved the exact unchanged WAV bytes.
- Cancellation, queued recovery after an actual process restart, saved clone,
  completed audio/hash, draft controls and dictionary persistence passed.
- MP3/M4A/WebM references decoded and created states; silence rejected. WAV
  reference validation and consent metadata passed. Speed/pitch extremes passed.
- Two-account private sharing/grant/use/revoke, owner-only reference access,
  output ownership, model-management role checks, anonymous/CSRF/traversal
  rejection and bounded/idempotent generation requests passed.
- Latest-schema restart and schema1-to2 migration regression passed; newer
  schemas are rejected without downgrading data. Replacement expansion is bounded.

## Managed acceptance and delivery status

The ordinary complete organization candidate contains36 apps/427 owned packages.
Actual native fresh selection and omission passed; omission retains only global
legal notices. MyAccounts OIDC, shared picker, restricted nested Pocket and MOSS
workers, genuine cloning, WAV/MP3, exact My Drive export, Markdown import and
Workspace global sign-out passed. A normal Voice-only candidate update retained
Workspace, Accounts and Drive bytes. Public110 rejected integration7 before
mutation; the new controller upgraded110→111, independently verified every
recovery file, rolled back and reapplied. Old app bytes, identities, paths,
ordinary files and Voice private data remained intact. Optional new Voice is an
explicit choice and is never automatically selected; setup exclusions persist.

Testing caught and fixed second-launch migration order, optional-app discovery,
and new-app rollback removing the app marker while keeping its user data. No
owner installation was a feature fixture. Actual publication/channel and both
owner-PC installation receipts remain pending at this source-tag checkpoint.
See subsequent repository delivery records for those separate stages.

## Limits

English is the native language acceptance coverage. The UI offers only the19
codes explicitly listed by MOSS; other languages need their own review. No Hindi,
noise suppression, GPU tier or network streaming is shipped. MOSS used about
1.7GiB tree working set and was slower than real time on this host; Pocket is the
lightweight preset path. Old4–8GB laptop performance is untested.

Voice quality, pronunciation and identity fidelity depend on the reference and
text. Signal checks are not a blind human listening study or proof of identity
match. No claim of quality parity/superiority to ElevenLabs or Pocket cloning is
made. Generated consented QA samples remain private for listening comparison.
The executable is unsigned. No service/reboot/clean-host acceptance is claimed.
