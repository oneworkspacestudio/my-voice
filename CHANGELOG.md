# Changes

## 0.1.0-preview.7

Choose saved or shared clones directly in the Speech voice dropdown. Selecting a
voice selects its matching model and preserves supported language choices.
Compatible model changes keep your chosen voice; switching engines uses the
matching preset. Segment voice choices stay compatible with the script model.

Automatic preview titles follow the selected voice, including older drafts with
a stale Alba title. Custom titles and historical recording data are preserved.
Each recording displays the voice(s) it actually used; changing the editor voice
does not relabel existing audio as the new speaker.

## 0.1.0-preview.6

Voice references are prepared automatically before cloning: convert the format,
trim silence and long pauses, balance volume safely and choose a useful excerpt
up to30 seconds. The original file is retained. The cloning screen explains the
changes and plays the prepared reference separately. Up to the first two minutes
of an uploaded file are analysed.

Quiet or longer usable recordings can continue without being recorded again.
Silence, insufficient speech and severe clipping now give specific guidance.
Preparation does not invent missing speech or reconstruct clipped peaks. An
existing browser reference can be retried after the update without reselecting it.

## 0.1.0-preview.5

Voice Cloning now shows installable models before they are installed. Install the
open MOSS model on this screen, watch download and preparation progress, and retry
setup without losing your reference or voice details. Recording and reference
upload unlock only when the selected cloning model is ready. Model readiness is
checked again before opening the microphone, accepting an upload or creating a
voice. The model list can be refreshed here; non-owners see owner-install guidance.

## 0.1.0-preview.4

Generate speech now uses the complete current script after editing or reopening
a recording. Stale chunks cannot silently replace new text. Voice previews play
separately and keep your draft, voice settings and recording intact.

Segment details are hidden by default. Use **Show segment details** to inspect or
edit them; hiding the panel keeps any intentional segment edits and voices.
Editing the main script returns generation to the complete text. Optional targeted
regeneration still retains the unchanged audio from other segments.

## 0.1.0-preview.3

Pocket speech now falls back to a native baseline x86-64 worker when the
upstream optimized Windows worker exits with an unsupported CPU instruction.
The same pinned model and voice files are reused. Windows worker failures show
the executable and numeric status with a relevant recovery message; private
speech, references and worker output are never included in those messages.

Models, saved voices, projects and history retain their existing format.

## 0.1.0-preview.2

First Windows preview: original Workspace interface inspired by familiar audio
tools, genuine CPU-only Pocket preset speech and open MOSS voice cloning,
sentence-level generation, pronunciation replacement, dialogue assignment,
targeted regeneration, TXT/Markdown import, batches, WAV/MP3/ZIP export, private
voice sharing, history, model management, OneAccount and Drive integration.

Requires Workspace preview111 for managed delivery. Optional noise suppression,
GPU tiers, Hindi, dubbing, voice agents and guaranteed voice-quality parity are
not included. See the native test report for measured hardware and limitations.

The earlier preview1 tag is source-only. This first binary release adds the full
Pocket Cargo.lock notice inventory and an accompanying source archive, including
its MPL-2.0 audio components. No published binary asset was replaced.
