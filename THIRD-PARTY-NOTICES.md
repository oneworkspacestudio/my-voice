# Third-party notices

My Voice source is MIT. Included support components retain their own licences.
No model weights or personal reference recordings are bundled in the executable.

| Included component | Version | Licence / notice |
|---|---|---|
| Go runtime and server modules | Go1.27.3; exact modules in docs/go-dependencies.json and go.sum | BSD/MIT/Apache notices in licenses/go and licenses/Go-BSD.txt |
| One Workspace apphost | v0.3.0-preview.110 | MIT; oneworkspacestudio/one-workspace |
| React / React DOM / scheduler | 19.3.1 / 19.3.1 / 0.28.0 | MIT; licenses/web |
| Lucide | 0.577.0 | ISC; licenses/web |
| Codedia AUI compiler | 0.2.1, build tool | MIT; licenses/web |
| Pocket native Rust CPU runtime | official v4.0.1 archive, CLI3.3.0 | MIT; licenses/Pocket-MIT.txt |
| MOSS CPU runtime source | 8b7bcc9341b3b4ef3a3a58ba1338a7d85ff133eb | Apache-2.0; licenses/MOSS-Apache-2.0.txt |
| CPython embedded x64 | 3.12.10 | PSF and incorporated notices; licenses/Python-3.12.txt |
| NumPy | 2.4.4 | BSD and incorporated numerical-library notices; licenses/python |
| ONNX Runtime | 1.24.4, CPU only | MIT and incorporated notices; licenses/python/onnxruntime |
| SentencePiece | 0.2.1 | Apache-2.0 |
| Flatbuffers / packaging / protobuf / SymPy / mpmath | exact pinned wheels in runtime.lock.json | Apache-2.0 / BSD / BSD / BSD / BSD, upstream wheel notices |
| FFmpeg audio runtime | 9.0.2, commit1e65c6cd5d | LGPL-2.1-or-later; licenses/FFmpeg-LGPL-2.1.txt |
| LAME | 3.100, statically included in FFmpeg | LGPL-2.0-or-later; licenses/LAME-LGPL-2.0.txt |

The small FFmpeg executable is built natively with Zig0.15.2, without GPL or
nonfree configuration, network protocols or autodetected external codecs.
Its complete FFmpeg/LAME corresponding sources and exact build recipe are offered
as `audio-runtime-corresponding-source.zip` beside every Windows release. The
recipe is also scripts/build-audio-runtime.sh; source/archive digests are locked
in runtime.lock.json. Users may modify/rebuild this LGPL program and replace the
`ffmpeg.exe` member in the source build's payload; prepare-runtime.py and the Go
source permit relinking/rebuilding the launcher. No restriction on reverse
engineering these libraries for debugging modifications is imposed.

Python support uses prebuilt pinned wheels without PyTorch. The pinned wheel
archives remain available from the URLs in runtime.lock.json. Each wheel's own
licences and bundled dependency notices apply. Models and the CC0 Open Voice
reference download separately on explicit acceptance: see MODEL-LICENSES.md.

No ElevenLabs code, branding, credentials or proprietary models are included.


Pocket's runtime incorporates Rust dependencies, including Symphonia audio
components under MPL-2.0. `licenses/Pocket-THIRD-PARTY-NOTICES.txt` retains the
upstream copyright/licence texts; `docs/pocket-dependencies.json` records the
complete pinned upstream Cargo.lock inventory (including optional/platform
crates). Not every locked crate is compiled into the Windows executable.
`pocket-runtime-corresponding-source.zip` supplies the runtime source and locked
dependency source beside this release. Test audio and model files are excluded;
the original upstream URLs and archive digests are included. This is in addition
to the separate LGPL FFmpeg/LAME source archive. The MPL-covered files remain
available under their original licence; the enclosing My Voice code stays MIT.
