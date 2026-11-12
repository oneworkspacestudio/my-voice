# Native Windows test report

## Preview7 direct cloned-voice selection

Public6 reproduces the missing clone: a real saved MOSS clone exists, but the
Alba/Pocket Speech dropdown offers only Alba. Exact7 retains the same clone and
allows direct dropdown selection from Alba; model, language and automatic title
are updated together. The submitted script and every generated segment use the
selected clone ID, with the complete two-sentence text. Real MOSS generation is
8.679s WAV plus MP3; the player shows the selected clone as the recorded speaker.

Eighteen frontend regressions cover complete script handling, cloning readiness,
cross-engine clone/shared/preset selection, compatible model/language retention,
custom/legacy automatic titles and immutable recording labels. Full Go tests
(with MYVOICE_TEST_ROOT set to the isolated NTFS fixture), vet and AUI/TypeScript/
Vite build pass. Desktop and phone-width dropdown/title layout are reviewed with
no horizontal overflow. Only public CC0 synthetic accounts/audio are fixtures.
Runtime/model bytes, schema2 and116 owned destinations remain unchanged. Managed
update, channel promotion and owner-target delivery receipts follow publication.

## Preview6 automatic reference preparation

Public5 rejects a recoverable26-second,44.1kHz stereo upload containing quiet
public CC0 speech and silence with the reported combined quality error. After
starting exact6 on the same synthetic saved data, the existing5 browser page
retries the identical selected file without reload/reselection and creates a
real clone. Original bytes remain exact; prepared audio is mono24kHz,6.38s,
RMS0.0661 with no clipped samples. A50-second reference prepares a valid30-second
excerpt and real clone. Silent input returns a specific explanation.

Exact6 synthesizes new two-segment speech from the prepared saved clone:7.381s
WAV and MP3. Six native preparation regressions cover quiet/padded audio, pauses,
excerpt selection, DC drift, severe clipping/silence/short input and near-fullscale
unclipped audio at the three-second boundary; originals are checked byte-for-byte.
Full Go tests/vet,11 frontend regressions and final AUI/TS/Vite build pass.
Native reference format conversion and real MOSS encoding are tested, without
owner recordings or physical microphone capture as fixtures. Runtime/model bytes,
schema2 and116 managed owned files are unchanged. Delivery checks follow release.

## Preview5 cloning setup before capture

Fresh native Windows model-manager download of763,682,544 pinned public MOSS
bytes completes in89.9s on this host, including real reference preparation.
An uninstalled model is visible in Voice Cloning with inline installation;
recording/upload remain disabled through download and preparation. A synthetic
network failure shows a retry action. Removing only the synthetic readiness
marker locks capture/creation while preserving the chosen reference and name;
restoring readiness unlocks them. A public CC0 ten-second reference is cloned
through the UI. The exact final executable retains that clone and generates
new two-segment speech:19.8605s WAV plus MP3 on CPU. No physical microphone or
owner reference/installation was a feature fixture.

Eleven frontend regressions, full Go tests/vet and final AUI/TypeScript/Vite build
pass. Native API checks preserve preinstall capability discovery and block
concurrent install/removal during preparation. Desktop and phone-width controls
are reviewed without horizontal overflow. A temporary Windows text-encoding
build error was corrected before final compilation; no failed artifact shipped.
Runtime/model bytes, database schema and116-package managed contract are unchanged.
Publication/managed-update/owner-target receipts follow separately after release.

## Preview4 full-script and preview isolation repair

The public3 bug was reproduced with synthetic data: preview replaced the draft,
then replacing the main text still generated the preview's cached chunks. The
same saved data reopened in4 with its full editor text and details collapsed.
Preview4 retained the draft during real voice preview, generated all3 sentences,
and generated a15-segment,75.320s recording through the native CPU worker. Every
chunk's interior samples, including the final chunk, matched the combined WAV;
total frames equalled every segment plus pauses. MP3 exports were also present.

Optional targeted regeneration retained all3 chunks and exact unchanged segment
bytes. Replacing the main editor after segment edits generated all4 new sentences
(9.589s audio). Six frontend regressions, native script-mode/punctuation tests,
the full Go suite/vet and AUI/TypeScript/Vite build pass. Desktop and a phone-sized
viewport (~390px) were visually reviewed without horizontal overflow. These are
synthetic isolated tests; no owner's installation, text or voices were fixtures.

Runtime/model bytes, storage schema and managed package contract are unchanged.
Managed update and owner-target delivery are recorded separately after release.

## Preview3 worker compatibility repair

On the same native host described below, the exact preview3 executable
`67f8aebe04f60dce52e48373844b8effbee9c9466b827865eec673b644e2adbf`
generated real WAV and MP3 for both Pocket18 and Pocket7, including 11 characters
without punctuation. The new baseline x86-64 worker also generated audible
24 kHz PCM in apphost's restricted native Job Object. A native process regression
simulates Windows STATUS_ILLEGAL_INSTRUCTION and verifies fallback, environment
scrubbing, ordinary error handling and suppression of private worker output.
All Go tests/vet and the unchanged production frontend build pass.

The deterministic runtime reconstruction preserves the complete preview2 payload
byte for byte when the portable executable is omitted. Rust1.99.1/MSVC native
compilation uses the same pinned source/Cargo.lock with baseline x86-64 and a
static CRT. The offline source recipe also passed. See runtime.lock.json and
the `workerRepair` record in docs/evidence/windows-native.json.

This host supports AVX2. These checks establish the baseline build and fallback
mechanism; they do not establish old physical CPU performance or the cause of a
remote user's generic preview2 failure. Managed rollout and target recovery are
recorded separately after publication, without publishing user text or audio.

## Preview2 baseline acceptance

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


## Preview8 model comparison acceptance (4 November 2026)

Nineteen frontend regressions, native Go tests/vet and production frontend build
pass. New coverage checks same-name model variants, engine/language isolation,
resumed and restarted downloads, interruption recovery, invalid ranges and
checksum rejection without promoting a partial file.

Both Qwen3-TTS Base Q8 checkpoints run with the pinned native Windows CPU runtime
inside apphost's supervised Job Object. Setup generates and checks real speech
before marking the model ready. Qwen1.7 was downloaded completely by the app from
its pinned public sources without authentication. Large files are streamed and
hashed without loading model weights into the Go heap. App-local Microsoft CRT
and OpenMP modules were observed loading from the verified embedded runtime.

A synthetic public7 MOSS voice survives in preview8. My Voices **Try another model**
reuses its owner-only saved reference; the original reference hash stays unchanged.
The UI creates separate 0.6B and1.7B variants with the same name, displays their
models in voice choices and recording captions, and preserves the full script.
Switching the own voice's Speech model opens the same reference-reuse flow.
Recording/upload stay disabled until the selected model is ready. Desktop and
390px presentation pass; no physical microphone was used.

| Model | Same two-sentence script: audio seconds | App wall seconds |
|---|---:|---:|
| Qwen3-TTS0.6B Q8 | 6.642 | 14.864 |
| Qwen3-TTS1.7B Q8 | 5.159 | 19.806 |

These include separate workers, both segments, postprocessing and WAV/MP3 exports,
with two CPU threads on this host. Download/setup time is excluded. RAM guidance
(4GB/8GB available) is an estimate, not a tested minimum specification. English
is tested; the other nine Qwen languages require their own listening review.
No blind similarity study, quality-superiority or low-memory/clean-host claim is
made. The larger model is an option to compare, not a guaranteed better match.

Managed public upgrade and each owner target are recorded separately in the
subsequent delivery receipt. No owner installation is a feature test fixture.

## Preview.9 Podcasts integration

An isolated NTFS fixture exercised MOSS and Pocket 18 with the packaged local
CPU runtime. Both stock narrator styles produced audible, non-silent 24 kHz
PCM. A signed-out Podcasts worker generated three complete scheduled WAVs and
resumed a previously reviewed script using the real local endpoint. Research
and writing were fixtures. Direct authorization, proxy denial and credential
revocation were verified. Nineteen browser regressions passed.
