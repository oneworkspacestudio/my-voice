"""Reconstruct embedded executable support assets from the committed digest lock.
Native development/build tooling only; ordinary app startup never fetches code.
Usage: python scripts/prepare-runtime.py I:/Projects/tooling/my-voice-runtime
"""
import hashlib,json,sys,urllib.request,zipfile,shutil
from pathlib import Path
repo=Path(__file__).resolve().parents[1];cache=Path(sys.argv[1]).resolve();cache.mkdir(parents=True,exist_ok=True)
lock=json.loads((repo/'runtime.lock.json').read_text(encoding='utf-8'))
def fetch(pin,name):
 p=cache/name
 if not p.exists():
  with urllib.request.urlopen(pin['url'],timeout=180) as src,p.open('wb') as dst:shutil.copyfileobj(src,dst,4<<20)
 h=hashlib.sha256(p.read_bytes()).hexdigest()
 if h!=pin['sha256']:raise RuntimeError('Digest mismatch: '+name)
 return p
stage=cache/'stage';stage.mkdir(exist_ok=True)
with zipfile.ZipFile(fetch(lock['pocket'],'pocket.zip')) as z:
 f=next(n for n in z.namelist() if n.endswith('/pocket-tts.exe'));(stage/'pocket-tts.exe').write_bytes(z.read(f))
with zipfile.ZipFile(fetch(lock['pocketPortable'],'pocket-portable-runtime.zip')) as z:
 (stage/'pocket-tts-portable.exe').write_bytes(z.read('pocket-tts-portable.exe'))
with zipfile.ZipFile(fetch(lock['ffmpeg'],'ffmpeg-runtime.zip')) as z:
 (stage/'ffmpeg.exe').write_bytes(z.read('ffmpeg.exe'))
# The native CPU build dispatches to a compatible CPU backend at runtime.
qwen=stage/'qwen';qwen.mkdir(exist_ok=True)
with zipfile.ZipFile(fetch(lock['qwenRuntime'],'qwen-runtime.zip')) as z:
 for name in z.namelist():
  base=Path(name).name
  needed=(base in ['llama-tts.exe','llama.dll','llama-common.dll','mtmd.dll','libomp.dll','LICENSE-LLVM-OpenMP'] or (base.startswith('ggml') and base.endswith('.dll') and 'rpc' not in base))
  if needed:(qwen/base).write_bytes(z.read(name))
# App-local release CRT from the pinned official Microsoft redistributable.
with zipfile.ZipFile(fetch(lock['qwenCRT'],'qwen-crt.vsix')) as z:
 for name in lock['qwenCRT']['members']:
  (qwen/name).write_bytes(z.read(lock['qwenCRT']['memberPrefix']+name))
runtime=stage/'moss';runtime.mkdir(exist_ok=True)
with zipfile.ZipFile(fetch(lock['python'],'python.zip')) as z:z.extractall(runtime)
(runtime/'python312._pth').write_text('python312.zip\n.\nLib/site-packages\nimport site\n')
site=runtime/'Lib/site-packages';site.mkdir(parents=True,exist_ok=True)
for wheel in lock['wheels']:
 with zipfile.ZipFile(fetch(wheel,wheel['name'])) as z:z.extractall(site)
source=cache/'moss-source.py'
url='https://raw.githubusercontent.com/OpenMOSS/MOSS-TTS-Nano/'+lock['mossSource']+'/ort_cpu_runtime.py'
if not source.exists():source.write_bytes(urllib.request.urlopen(url,timeout=60).read())
if hashlib.sha256(source.read_bytes()).hexdigest()!=lock['mossRuntimeSHA256']:raise RuntimeError('MOSS runtime source digest mismatch')
shutil.copyfile(source,runtime/'ort_cpu_runtime.py');shutil.copyfile(repo/'runtime-src/worker.py',runtime/'worker.py')
with zipfile.ZipFile(repo/'internal/app/payload.zip','w',zipfile.ZIP_DEFLATED,compresslevel=6) as z:
 for p in sorted(stage.rglob('*')):
  if not p.is_file() or '__pycache__' in p.parts:continue
  info=zipfile.ZipInfo(p.relative_to(stage).as_posix(),date_time=(1980,1,1,0,0,0));info.compress_type=zipfile.ZIP_DEFLATED;info.external_attr=0o600<<16;z.writestr(info,p.read_bytes())
print('Prepared verified embedded runtime. Models remain separate downloads.')
