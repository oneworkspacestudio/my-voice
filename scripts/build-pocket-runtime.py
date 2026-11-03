"""Build Pocket's pinned source for baseline x86-64 on native Windows.

Run in an x64 MSVC developer shell with Rust/Cargo and CMake on PATH:
  python scripts/build-pocket-runtime.py --sources <pocket-crates> --work C:/... --out <zip>
The source directory is the exact archive set in pocket-runtime-corresponding-source.zip.
No model weights, code downloads, Cargo.lock changes or global configuration writes.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import tarfile
import zipfile
try:
    import tomllib as toml
except ImportError:
    import tomli as toml

p = argparse.ArgumentParser(description=__doc__)
p.add_argument('--sources', type=Path, required=True)
p.add_argument('--work', type=Path, required=True)
p.add_argument('--out', type=Path, required=True)
p.add_argument('--cargo', default='cargo')
a = p.parse_args()
if os.name != 'nt':
    raise RuntimeError('This release recipe requires native Windows/MSVC')
repo = Path(__file__).resolve().parents[1]
audit = json.loads((repo / 'docs/pocket-dependencies.json').read_text())
source = a.sources.resolve()
work = a.work.resolve()
work.mkdir(parents=True, exist_ok=True)

def verified(path, digest):
    raw = path.read_bytes()
    if hashlib.sha256(raw).hexdigest() != digest:
        raise RuntimeError('Source digest mismatch: ' + path.name)
    return raw

def safe(name):
    name = name.replace('\\', '/')
    if name.startswith('/') or ':' in name or '..' in Path(name).parts:
        raise RuntimeError('Invalid archive member')

def unzip(path, dest):
    with zipfile.ZipFile(path) as z:
        for item in z.infolist():
            safe(item.filename)
            if not item.is_dir():
                target = dest / item.filename
                target.parent.mkdir(parents=True, exist_ok=True)
                target.write_bytes(z.read(item))

verified(source / 'pocket-source.zip', audit['sourceSHA256'])
unzip(source / 'pocket-source.zip', work / 'upstream')
upstream = work / 'upstream' / ('pocket-tts-' + audit['sourceCommit'])
project = upstream / 'rust'
lock_bytes = verified(project / 'Cargo.lock', audit['cargoLockSHA256'])
lock = toml.loads(lock_bytes.decode())
vendor = work / 'vendor'
vendor.mkdir(exist_ok=True)

def checksum(target, package):
    files = {f.relative_to(target).as_posix(): hashlib.sha256(f.read_bytes()).hexdigest()
             for f in target.rglob('*') if f.is_file() and f.name != '.cargo-checksum.json'}
    (target / '.cargo-checksum.json').write_text(json.dumps({'files': files, 'package': package}))

for package in lock['package']:
    if not package.get('source', '').startswith('registry+'):
        continue
    name = package['name'] + '-' + package['version']
    arc = source / (name + '.crate')
    verified(arc, package['checksum'])
    target = vendor / name
    if not (target / '.cargo-checksum.json').exists():
        with tarfile.open(arc) as t:
            for member in t.getmembers():
                safe(member.name)
                if member.issym() or member.islnk():
                    raise RuntimeError('Unexpected crate link')
            t.extractall(vendor)
        checksum(target, package['checksum'])

git_dep = next(x for x in audit['packages'] if x['name'] == 'ptts')
verified(source / 'ptts-git-source.zip', git_dep['sha256'])
unzip(source / 'ptts-git-source.zip', work / 'git')
git_root = next((work / 'git').iterdir())
workspace = toml.loads((git_root / 'Cargo.toml').read_text())['workspace']
target = vendor / ('ptts-' + git_dep['version'])
shutil.copytree(git_root / 'ptts', target, dirs_exist_ok=True)
# Cargo's directory source needs a package manifest. Resolve workspace inheritance
# exactly, retaining all dependency versions/features and unmodified library code.
raw = (target / 'Cargo.toml').read_text()
def value(v):
    if isinstance(v, dict):
        return '{ ' + ', '.join(k + ' = ' + value(x) for k, x in v.items()) + ' }'
    return json.dumps(v)
for key, val in workspace['package'].items():
    raw = raw.replace(key + '.workspace = true', key + ' = ' + value(val))
def dependency(match):
    name = match[1]
    local = toml.loads('v = ' + match[2])['v']
    base = workspace['dependencies'][name]
    base = {'version': base} if isinstance(base, str) else base.copy()
    local.pop('workspace')
    base.update(local)
    return name + ' = ' + value(base)
raw = re.sub(r'(?m)^([\w-]+) = (\{ workspace = true[^\n]+\})$', dependency, raw)
# This dev-only example dependency is outside the package; it is not built/linked.
raw = re.sub(r'(?m)^ptts-coreml = .*\n', '', raw)
(target / 'Cargo.toml').write_text(raw)
checksum(target, None)
git_package = next(x for x in lock['package'] if x.get('source', '').startswith('git+'))
git_url, revision = git_package['source'][4:].split('?rev=')
revision = revision.split('#')[0]
cfg = upstream / '.cargo'
cfg.mkdir(exist_ok=True)
(cfg / 'config.toml').write_text(
    '[source.crates-io]\nreplace-with="vendored-sources"\n'
    '[source.ptts-git]\ngit=' + json.dumps(git_url) + '\nrev=' + json.dumps(revision) +
    '\nreplace-with="vendored-sources"\n[source.vendored-sources]\ndirectory=' +
    json.dumps(vendor.as_posix()) + '\n')
env = dict(os.environ, RUSTFLAGS='-C target-cpu=x86-64 -C target-feature=+crt-static',
           CARGO_TARGET_DIR=str(work / 't'))
subprocess.run([a.cargo, 'build', '--release', '--locked', '--offline', '--target',
                'x86_64-pc-windows-msvc'], cwd=project, env=env, check=True)
binary = (work / 't/x86_64-pc-windows-msvc/release/pocket-tts.exe').read_bytes()
a.out.parent.mkdir(parents=True, exist_ok=True)
with zipfile.ZipFile(a.out, 'w', zipfile.ZIP_DEFLATED, compresslevel=6) as z:
    info = zipfile.ZipInfo('pocket-tts-portable.exe', date_time=(1980, 1, 1, 0, 0, 0))
    info.compress_type = zipfile.ZIP_DEFLATED
    info.external_attr = 0o600 << 16
    z.writestr(info, binary)
print(json.dumps({'sha256': hashlib.sha256(a.out.read_bytes()).hexdigest(),
                  'executableSHA256': hashlib.sha256(binary).hexdigest(),
                  'size': a.out.stat().st_size, 'sourceCommit': audit['sourceCommit'],
                  'rustflags': env['RUSTFLAGS']}))
