"""Maintenance: reconstruct Pocket dependency notices from its exact upstream release.
Usage: python scripts/collect-pocket-notices.py <TASK_CACHE> (Python3.11+ or tomli).
"""
import urllib.request,urllib.error,subprocess,json,hashlib,re,tarfile,io,zipfile,concurrent.futures,time,sys
from pathlib import Path
try:import tomllib as toml
except ImportError:import tomli as toml
art=Path(sys.argv[1]).resolve();art.mkdir(parents=True,exist_ok=True);cache=art/'pocket-crates';cache.mkdir(exist_ok=True);gh='gh'
commit=json.loads(subprocess.check_output([gh,'api','repos/kyutai-labs/pocket-tts/commits/v4.0.1']))['sha'];lock=urllib.request.urlopen('https://raw.githubusercontent.com/kyutai-labs/pocket-tts/'+commit+'/rust/Cargo.lock').read();(art/'pocket-Cargo.lock').write_bytes(lock);packages=toml.loads(lock.decode())['package'];results=[];texts=[]
def fetch(url,path,expected=None):
 if not path.exists():
  for n in range(3):
   try:
    raw=urllib.request.urlopen(url,timeout=120).read();path.write_bytes(raw);break
   except Exception:
    if n==2:raise
    time.sleep(.5)
 raw=path.read_bytes();sha=hashlib.sha256(raw).hexdigest()
 if expected and sha!=expected:raise RuntimeError('Upstream digest mismatch '+path.name)
 return raw,sha

def collect(p):
 name,version=p['name'],p['version'];source=p.get('source','');assert re.fullmatch('[A-Za-z0-9_-]+',name)
 if not source.startswith('registry+'):return None
 url='https://static.crates.io/crates/'+name+'/'+name+'-'+version+'.crate';raw,sha=fetch(url,cache/(name+'-'+version+'.crate'),p['checksum']);notices=[];meta={}
 with tarfile.open(fileobj=io.BytesIO(raw),mode='r:gz') as tar:
  for m in tar.getmembers():
   if not m.isfile():continue
   base=m.name.rsplit('/',1)[-1].upper()
   if m.name==name+'-'+version+'/Cargo.toml':meta=toml.loads(tar.extractfile(m).read().decode())
   if base.startswith(('LICENSE','LICENCE','COPYING','NOTICE','COPYRIGHT')):
    data=tar.extractfile(m).read()
    if len(data)<2_000_000:notices.append((m.name,data.decode('utf-8',errors='replace')))
 return ({'name':name,'version':version,'source':url,'sha256':sha,'license':meta.get('package',{}).get('license'),'noticeFiles':[n for n,t in notices]},notices)
print('Collecting exact Pocket runtime lock: '+str(len(packages))+' packages',flush=True)
with concurrent.futures.ThreadPoolExecutor(max_workers=8) as pool:
 for i,value in enumerate(pool.map(collect,packages)):
  if value:
   record,notices=value;results.append(record);texts.append('\n\n'+'='*70+'\n'+record['name']+' '+record['version']+'\nSPDX: '+str(record['license'])+'\nSource: '+record['source']+'\nSHA256: '+record['sha256']+'\n'+''.join('\n--- '+n+' ---\n'+t for n,t in notices))
  if i and i%50==0:print('Inventoried '+str(i)+' packages',flush=True)
# Pinned Git dependency source, in addition to registry crates.
gitdep=next(p for p in packages if p.get('source','').startswith('git+'));gitrev=gitdep['source'].split('#')[1];url='https://codeload.github.com/gabrieldemarmiesse/xn-ptts/zip/'+gitrev;raw,sha=fetch(url,cache/'ptts-git-source.zip');notices=[]
with zipfile.ZipFile(io.BytesIO(raw)) as z:
 for n in z.namelist():
  if n.rsplit('/',1)[-1].upper().startswith(('LICENSE','LICENCE','COPYING','NOTICE','COPYRIGHT')):notices.append((n,z.read(n).decode('utf-8',errors='replace')))
results.append({'name':'ptts','version':gitdep['version'],'source':url,'sha256':sha,'commit':gitrev,'license':'see source notices','noticeFiles':[n for n,t in notices]});texts.append('\n\nptts '+gitrev+'\n'+''.join('\n--- '+n+' ---\n'+t for n,t in notices))
pocketurl='https://codeload.github.com/kyutai-labs/pocket-tts/zip/'+commit;raw,sha=fetch(pocketurl,cache/'pocket-source.zip')
manifest={'schema':1,'upstreamRelease':'v4.0.1','runtimeVersion':'3.3.0','sourceCommit':commit,'sourceURL':pocketurl,'sourceSHA256':sha,'cargoLockSHA256':hashlib.sha256(lock).hexdigest(),'scope':'Full upstream Cargo.lock including optional/platform-specific dependencies; not every crate is in the Windows binary','packages':results}
(art/'pocket-dependencies.json').write_text(json.dumps(manifest,indent=2)+'\n',encoding='utf-8');(art/'Pocket-THIRD-PARTY-NOTICES.txt').write_text('Pocket TTS native runtime dependency notices\nFull locked source provenance: pocket-dependencies.json\n'+''.join(texts),encoding='utf-8');print('Pocket dependency audit complete '+str(len(results))+' entries',flush=True)
print('Licences',sorted({str(p.get('license')) for p in results}),flush=True)
