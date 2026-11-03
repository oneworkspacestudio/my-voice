"""Package reviewed source documents and one tested executable, never private data."""
import argparse, hashlib, json, subprocess, zipfile
from pathlib import Path

root = Path(__file__).resolve().parents[1]
p = argparse.ArgumentParser()
p.add_argument('--binary', required=True, type=Path)
p.add_argument('--out', required=True, type=Path)
p.add_argument('--candidate', action='store_true')
a = p.parse_args()
git = ['git', '-c', 'safe.directory=' + root.as_posix(), '-C', str(root)]
contract = json.loads((root / 'compatibility.json').read_text(encoding='utf-8-sig'))
if a.candidate:
    commit = 'unpublished-candidate'
else:
    if subprocess.check_output(git + ['status', '--porcelain']).strip():
        raise RuntimeError('Commit the reviewed source before packaging a release')
    commit = subprocess.check_output(git + ['rev-parse', 'HEAD'], text=True).strip()
binary = a.binary.read_bytes()
if len(binary) < 1_000_000 or binary[:2] != b'MZ':
    raise RuntimeError('Expected the tested native Windows executable')
names = ['compatibility.json', 'LICENSE', 'NOTICE', 'THIRD-PARTY-NOTICES.md',
         'README.md', 'INSTALL-WINDOWS.md', 'MODEL-LICENSES.md',
         'TEST-REPORT-WINDOWS.md', 'ARCHITECTURE.md', 'SECURITY.md', 'CHANGELOG.md',
         'docs/API.md', 'docs/go-dependencies.json', 'docs/web-dependencies.json',
         'docs/evidence/windows-native.json', 'docs/pocket-dependencies.json']
entries = {'MyVoice.exe': binary}
for name in names:
    if not (root / name).is_file():
        raise RuntimeError('Missing required package document: ' + name)
for path in [root / name for name in names] + sorted((root / 'licenses').rglob('*')):
    if not path.is_file():
        continue
    if path.is_symlink() or not path.resolve().is_relative_to(root):
        raise RuntimeError('Unexpected external package input')
    entries[path.relative_to(root).as_posix()] = path.read_bytes()
manifest = {'schema': 1, 'version': contract['version'], 'sourceCommit': commit,
            'candidate': a.candidate, 'builtNatively': True,
            'executable': 'MyVoice.exe', 'modelWeightsBundled': False,
            'executableSHA256': hashlib.sha256(binary).hexdigest(),
            'runtimePayloadSHA256': hashlib.sha256((root / 'internal/app/payload.zip').read_bytes()).hexdigest()}
entries['manifest.json'] = (json.dumps(manifest, indent=2) + '\n').encode()
a.out.parent.mkdir(parents=True, exist_ok=True)
with zipfile.ZipFile(a.out, 'w', zipfile.ZIP_DEFLATED, compresslevel=6) as z:
    for name, content in sorted(entries.items()):
        info = zipfile.ZipInfo(name, date_time=(1980, 1, 1, 0, 0, 0))
        info.compress_type = zipfile.ZIP_DEFLATED
        info.external_attr = 0o600 << 16
        z.writestr(info, content)
print(json.dumps({'file': a.out.name, 'size': a.out.stat().st_size,
                  'sha256': hashlib.sha256(a.out.read_bytes()).hexdigest(),
                  'candidate': a.candidate}))
