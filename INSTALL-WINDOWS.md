# Windows installation

## Managed One Workspace

Select **My Voice** under Creative & Development during Workspace setup. To add
it to an existing installation, update Workspace to preview112 or newer, then
choose My Voice in Updates. Open it from the Workspace Apps picker. Existing
MyAccounts sign-in, Workspace origin and protected storage come from the manager.
Already installed Voice2 remains compatible with Workspace111.

Go to Settings, read the model licence, and choose Download. Use MOSS for private
voice cloning and multilingual generation. Pocket 18.2M/6.9M are smaller English
preset options. Downloaded models can generate offline; managed sign-in still
requires the configured identity service when a session must reconnect.

## Standalone portable package

Extract the download to a user-controlled program directory and run MyVoice.exe.
It listens on loopback and opens the browser. Its one-time setup link creates the
first local account. Passwords use bcrypt; the local account is not a OneAccount.
The normal data location is `%LOCALAPPDATA%/OneWorkspace/MyVoice` on NTFS.
It is separate from the replaceable EXE and contains private voice/audio data.

Advanced example (create a new dedicated NTFS folder):

```powershell
./MyVoice.exe --data C:/MyVoicePrivate --listen 127.0.0.1:9038 --no-browser
```

Do not bind the app directly to a LAN address. Use the authenticated Workspace
proxy for another browser/device. Standalone Drive integration is unavailable.
The executable is unsigned; no production signing certificate is included.

## Update, backup and removal

Managed updates use Workspace's normal protected update/recovery flow and retain
unselected apps, installation identity and private libraries. Backups are a
separate encrypted Workspace operation. Standalone users can use apphost's
`--backup`/`--restore` with `WORKSPACE_BACKUP_PASSWORD` (never put it in source).

Remove a model in Settings to reclaim its download space. Remove individual
recordings from History after export. Deleting a voice removes its references,
model-specific state and grants, while existing generated recordings are retained.
Removing the program does not delete private data. Preserve your data directory
before an uninstall. No service, driver, firewall exemption or reboot is required.

## Build from source on Windows

Install Node and Go1.27.3+, using a task-owned build cache. Run
`python scripts/prepare-runtime.py <CACHE>` to reconstruct digest-locked support
assets, then `scripts/build.ps1 -Go <GO_EXE> -Out <DESTINATION>`. This does not
download model weights. Use `MYVOICE_TEST_ROOT` and `GOTMPDIR` on dedicated NTFS
test directories for native tests; never use an owner's data as a fixture.
