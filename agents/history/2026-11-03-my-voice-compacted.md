# My Voice implementation and delivery

User authorizes the supplied ElevenLabs-inspired My Voice brief/UI, Workspace
design system, source push/release and installation on both PCs. User accepted
Pocket web gate but local token returned401/browser download blocked. They then
explicitly requested an open alternative. Use MOSS ONNX for cloning, ungated
Pocket18/7 presets for lightweight English. Never mirror gated weights or claim
quality parity/superiority. No agents or Actions were requested/used.

New actual app root I:/Projects/my-voice, isolated suite worktree
I:/Projects/workspace-my-voice based on public59ea61d/main. Aggregate outer repo
must never be published. SDK is pinned to public Workspace110; no local replace.
Suite integration targets111/integration7/voice port base+38, optional Creative
category, apphost OIDC/storage/recovery/update ownership, Apps picker and AI
native connector. Source only at this checkpoint; no push/release/owner install.

Implemented AUI/React UI, native Go apphost, CPU engines, verified model manager,
owner-scoped SQLite jobs/projects/voices, permission/consent/grants, WAV/MP3/ZIP,
sentence chunks, batch/cancel/retry/resume, targeted reuse of unchanged bytes,
snapshot dictionaries/project overrides, request idempotency, private Drive API
subject validation/import/export reconciliation. Optional noise cleanup/GPU/Hindi
are absent. MOSS upstream table names19 codes despite advertising20; expose only
those codes, native English coverage only. Runtimes embedded; models separate.

Native source/unit/CSRF/scope/queue/path checks and packaged CPU speech/cloning
pass using isolated NTFS C:/Projects/.cache/my-voice-20261102, artifacts on
I:/Projects/tooling/my-voice-20261102. Original heavy-runtime candidate is still
running there for acceptance; native-process.json tracks exact PID. CC0 donated
references only; owner installation is never a feature fixture. Standalone two
account private sharing/revoke and actual distinct-reference dialogue pass.
Targeted regeneration preserves exact unchanged WAV segment bytes.

First long-text run hit an intermittent worker/OS failure; bounded clean worker
retries plus sanitized OS-status logging added, and long text/ZIP later passed.
One cancellation test lost its HTTP response; do not call it passed yet. New
idempotent enqueue supports safe request retries. MOSS first postprocessing run
incorrectly halved48k sample rate; fixed with resampling before rate/pitch filters.
Two atempo stages now cover all allowed speed/pitch combinations. Final exact
binary acceptance and restart/cancel/audio-format checks remain required.

Desktop1600/phone390 packaged UI flows passed. Visual QA caught labels on icon
buttons and a mid-transition phone screenshot; fixed CSS and captures must be
rerun. Tree-shaken frontend is now~280KB JS. Current final source has subsequent
draft snapshot, Drive import and migration2 changes requiring final rebuild/tests.

Suite tests initially failed because Go1.27 makes testing.TempDir obey GOTMPDIR;
TEMP/TMP alone do not move it from exFAT. Set GOTMPDIR to synthetic C NTFS for
tests, keep general large builds on I. catalog/distribution/suite/branding/apphost
then pass. Voice storage-profile registration was corrected by the registry gate.

Native minimal FFmpeg/LAME build is in progress under task-owned C codec root,
using Git Bash (not WSL), MSYS make and portable Zig. Script/log under tooling
codec; retain corresponding source and recipes. It replaces the provisional large
BtbN payload for size/licensing clarity; no VM/container/tool system install.
Update runtime lock/prepare recipe/notices and deterministic payload once it passes.

Next: finish minimal audio build, exact final native/UI/performance/restart and
formats; document honest measurements; app immutable preview1 publication and
anonymous digests; pin to suite normal rebuild path; synthetic actual installer
fresh/110 upgrade/rollback with real Accounts+Drive; mandatory signed Android12
publisher for suite111; channel verification; normal protected installs on both
PCs and per-target served byte/health/recovery receipts. Primary is Workspace110
under C:/tmp/workspace/Private, public workspace.m8e.co.uk,35 app contracts. Second
PC is workspace-2.m8e.co.uk. Preserve all owner state/credentials/unselected apps
and newest verified recovery. Use installed updater repair handoff when old
integration6 needs new installer7; never edit installed catalog manually.


## 3 November native gates and managed-candidate checkpoint

No commits/push/releases/channel promotion/owner Voice installation yet. New public
main e465a8a was fast-forwarded into suite worktree: preserve MyAI20 lock/channel
and both PCs' independently installed MyAI20. The public110 full bundle embeds
MyAI14/Social6; existing repair/update code preserves newer installed app bytes.
Secondary baseline110/MyAI20/Hosting2, pre-existing MyMeet3 “Review release” warning.

MOSS now uses only explicit locally saved state or the CC0 Open Voice preset,
reference pinned to kyutai/tts-voices323332d.../voice-donations/0a67.wav. No upstream
named demo state is offered. MOSS install verifies that reference then prepares
preset-state.json. Language metadata/filter added. Models validate marker version
and asset presence/size; gated Pocket cloning remains absent.

Native minimal FFmpeg succeeded:3,202,048 byte executable,1,510,883 byte ZIP.
FFmpeg1e65c6cd + LAME3.100, LGPL source/recipe/notices retained. Runtime ZIP now
54,349,265 bytes, not107MB. runtime.lock points to intended future immutable own
release ffmpeg-runtime.zip; MUST publish that asset and complete source ZIP before
calling source reconstruction public-ready. prepare-runtime.py succeeded. Raw
upstream MOSS runtime digest99477efa... differs from Windows CRLF checkout digest;
lock corrected to raw unmodified upstream bytes. A cp1252 UI insertion was repaired.

Restart caught a real migration-order bug: SDK Migrate(1) rejected schema2 at
second launch. New startup reads current version, applies only pending migrations,
rejects newer schema. Regression covers1→2/restart/preserved project/future reject.
Native actual process restart now preserves clone, audio SHA, draft controls and
resumes interrupted segments with exact completed byte reuse. Windows test harness
PIPE output was inherited by background processes, causing an EOF wait despite
successful restart; use file stdout/stderr and bounded process wait, not PIPE.
restart-only.py passes. MP3/M4A/WebM state creation/deletion, WAV/silence validation,
speed/pitch extremes, dictionaries/dialogue/regeneration/long text/batch/cancel,
sharing/revoke, permissions, desktop/phone UI pass. TryLock rejects concurrent
reference preparation before audio read; replacement expansion bounded in API/UI.
Empty existing draft autosave corrected; final UI production build passes.

All9 CPU benchmark cases complete, two threads. Pocket18:13.012s audio/2.280s wall,
92.055/16.253 long,~464MiB peak tree. Pocket7:13.252/2.255,87.255/14.888 long,
~254MiB. MOSS:28.444/39.170,13.328/22.689 second request,146.809/292.786 long,
~1.7GiB. No network first-audio exposure; playback after assembly; no low-end,
physical-mic/phone, human blind listening or quality-superiority claim. Exact
candidate benchmark SHA/runtime provenance retained privately and sanitized in
app docs/evidence/windows-native.json. Startup to health existing data~1.6s.

Latest app binary tooling/my-voice-20261102/MyVoice-candidate-final.exe includes
all source fixes; MyVoice-ship.exe is running fixture (diff only empty-draft UI fix).
Candidate package voice-candidate.zip created with complete notices/docs and
manifest sourceCommit=unpublished-candidate. TEST-REPORT-WINDOWS.md explicitly
marks managed/publication/owner delivery pending. No real user audio in repo.

Important: suite organization-components.lock.json myVoice currently contains a
LOCAL FUTURE candidate pin. Do not commit/push it until public asset verified.
Added local-only --voice-candidate / -MyVoiceCandidate through ordinary rebuild
path; managed-build.candidate=true makes publisher reject such outputs. This allows
actual managed OIDC/Drive/fresh/upgrade/rollback acceptance BEFORE app release.
Final normal rebuild must use verified public pin without candidate flag. Voice
catalog port corrected9028 (base+38), standalone9038 remains separate. Android12
now includes generated brand_my_voice.png. Four unrelated generated exports restored.

First candidate rebuild rejected a partial old seed cache348MB, correctly.
Official public103 expected1,215,056,385 bytes, pinned cc61e4bf...; full anonymous
download session68040 underway to suite-cache/*.download, then SHA verified/rename.
Old invalid copy retained *.invalid-cached, original other-task cache untouched.
Rerun scripts/rebuild.ps1 -MyVoiceCandidate ...voice-candidate.zip with cache/out
under tooling root. Candidate builder uses public baseline103 + public remote10,
Tickets1, Authenticator2, Accounts, Hosting2; add Voice. Final complete inventory36.
Next actual synthetic setup select/omit Voice (Accounts/Drive/Voice), real SSO and
Drive export/import, protected public110→111 upgrade/rollback/reapply. Retain exact
unselected bytes/identity/paths. Then app source+annotatedtag/release (all hashes,
corresponding sources/runtime ZIP), replace local pin with public digest, normal
rebuild, mandatory signed Android publisher111/channel, primary normal new-installer
protected update and secondary owner Updates repair→optional Voice selection.

CUA Chrome3 initialized, documentation rewritten. New ephemeral task tab585734847
voiceRolloutTab currently secondary /updates, authenticated; voiceSecondaryBaseline
stores34 table rows for later comparison. MyMeet review warning pre-existed.
Original HF blocked-download tab still task-modified; restore original model tree
only if still that blocked page, avoid user work. No owner model download/consent
has been performed. Models remain explicit first-use owner installation.


## Later checkpoint: managed acceptance, optional upgrade repair, lazy CPU runtime

Fresh candidate actual installer passes Accounts OIDC/owner binding/anonymous deny,
shared Apps picker, managed Pocket download+CPU WAV/MP3, exact My Drive export and
Markdown import, global logout. Extended managed test also installs open MOSS,
creates CC0 clone, synthesizes distinct text and deletes its state in the nested
restricted worker. managed-acceptance.json passes; all fixture services stopped.
Initial Apps/Drive waits were test locator assumptions (manager converts Apps link
to button; success contains destination), not app defects. Omission also passes:
no Voice executable/runtime/contract/docs/private data. Suite intentionally keeps
all legal licence notices globally; do not assert those are absent.

Actual public110 controller rejects incoming integration7 before config mutation.
New111 candidate then exposed platform gap: optional new apps were never offered
by markNewDefaults/WantedApps, so explicit voice selection said unavailable.
Fixed distribution UpdateChoice.Optional, new optional-group discovery respecting
explicit exclusions/dependency exclusions; optional groups never auto-selected by
CLI, Workspace default addition or startup bridge. Maintenance UI leaves optional
checkbox unchecked and enabled alongside Workspace, with explicit hint. Source
helper renamed markNewApps. Existing coupled-group selections are preserved.
New regression covers explicit Voice/namespace/rebased port, default omission and
opt-out retention. distribution/maintenance/catalog/suite/branding all now pass.
Corrected test-only expected29028→29038; actual catalog Voice port9028/base+38.

Cold benchmark found unused MOSS runtime verification on Pocket first request
(~25.5s). App runtimeFor(false) now prepares/verifies only Pocket+FFmpeg; MOSS calls
runtimeFor(true), cached separately. Exact final app binary now MyVoice-lazy-candidate.exe
SHA0913a566d5267060cb485d8554f49ccccdb952e82c7ab6677ccb7f75c91f1f87.
Payload SHA87996d2929a2db9f3aee664e7b769a473315638de1426592bc6b58c9ef93c76e unchanged.
Fresh-runtime native proof confirms no moss/ directory after Pocket; explicit MOSS
then prepares verified Python and synthesizes (4.1555s audio,19.969s first-use wall).
Native Pocket final workloads329 chars/~19s and1963 chars/~2min now measured:
18M cold19s/4.621s; long22.174s;7M long119.129s/20.002s. Full results in
benchmarks-pocket-final.json; before-lazy evidence retained separately. Need update
public report/evidence to final metrics. Phone library/settings overflow checks pass.
Private4s CC0 sample open-moss-sample.mp3 shown to user for listening, no superiority claim.

New source scripts/package-windows.py whitelists docs/licences+tested EXE, writes
fixed-timestamp ZIP and provenance, public mode requires clean committed source.
voice_component now validates manifest/native SHA and rejects candidate/noncommit
metadata on normal public path; allow_candidate only explicit development flag.
Current voice-candidate.zip59,821,840 bytes SHAae6d7499e46373d57a553d13550bd88cec058d07a84a1771e526d229b0266276
includes lazy binary, candidate metadata. Local FUTURE myVoice lock pin updated;
still MUST NOT push/commit until actual public archive exists and digest verified.
Suite candidate rebuild session87285 running with latest optional fix and newapp.
After complete run upgrade-acceptance.py --resume: existing pristine public110
fixture C:/Projects/.cache/my-voice-20261102/upgrade-1793668698752/Private,
programs tooling/.../upgrade-programs-1793668698752, origin localhost40020.
Script now supports resume; no owner fixtures. It will test protected110→111,
all old payload/identity/path retention, independent backup files, rollback,
model/project retention and reapply. Prior failure was before mutation; root110.
Then use normal Voice-only candidate update on the existing managed111 fixture to
exercise independent app update and lazy nested runtime, rerun managed acceptance.
Final public-source/channel/owner deployment work remains unchanged/pending.
Primary inspected110/MyAI20/35apps/noVoice, completed journal; secondary same110,
MyAI20, pre-existing Meet review warning. No owner deployment/model consent yet.


## Recovery reapply bug fixed before publication

Fresh managed test after lazy Voice-only update passes again, including MOSS and
Drive (the harness now binds the exact POST-created job ID; reading latest history
raced a prior MOSS result, so the earlier export mismatch was a test bug).

Actual protected110→111 upgrade and rollback passed, but reapply revealed Voice
failed while the updater still reported portal ready. Rollback retained models/
projects yet removed new app.json/DB, leaving a non-empty folder which apphost
correctly refused to adopt. Do not bypass that guard or patch markers manually.
Source now derives exclusively new private roots from the two verified recovery
snapshots and preserves the complete root only when absent from the old catalog/
snapshot and unshared with old apps. Existing app schema recovery is unchanged.
Regression tests preserve new Voice marker/DB/projects, still clear old metadata,
and reject preservation of pre-existing/snapshot-backed roots. Added explicit
Voice /healthz readiness to update restart (like Tickets/Forms), so this selected
managed app cannot fail silently behind a healthy portal. All suite/distribution/
maintenance/catalog/branding tests pass. No owner install affected.

Old failed upgrade fixture upgrade-1793668698752 is stopped and intentionally not
reused/repaired. New candidate rebuild+fresh full110 setup/upgrade/rollback/reapply
runs in exec session56532, logs suite-candidate-build.log then upgrade-*.log;
upgrade-acceptance.py without --resume creates a new synthetic root and overwrites
upgrade-process.json. Await it. App bytes remain exact0913a566...; only suite code
changed. Re-run managed smoke only if resulting code changes justify it; current
app-only proof and managed MOSS/Drive are green. Still no public repo/push/release
or owner Voice installs. Suite public111 remains absent, main e465a8a unchanged.


Managed candidate upgrade/rollback/reapply now fully passes on a fresh native110 fixture. All old app payloads/identity/paths and new Voice private state are retained. New app marker/database preservation and explicit Voice readiness are verified. Final app remains0913a566...; no further app runtime code changes. Publication/owner installs pending. Full suite tests/vet and publisher regression underway.

All full-suite Go tests, touched-package vet, app vet, publisher regression, branding desktop/phone and image checks pass. Public my-voice repository created with Actions explicitly disabled. Preparing initial source commit/tag and immutable app assets; both owner installations remain pending. Copied upstream licence whitespace is retained verbatim; own-source whitespace checks are clean.

Final preview2 native smoke and desktop/phone UI pass. Exact EXE SHA256 e66d61720002393890a8ab46273919a04cb95dbf716f217f88ae84f6a23b4189. Official Pocket transitive MPL components now have Cargo.lock-derived dependency notices and corresponding source offer. Source-only preview1 remains unchanged; first downloadable release is preview2. Runtime payload unchanged. Publication and both production installs pending; managed candidate acceptance already passes.
