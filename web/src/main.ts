import React,{useEffect,useRef,useState} from 'react';

import {createRoot} from 'react-dom/client';

import {Voice} from './generated/Voice';

import {Bindings} from './components/ui';
import {buildScript,editedSegments,loadDraft,segmentFields,splitEditorText,scriptSignature} from './editor';

import {cloningSetup,requireCloningModel,modelInstalling} from './cloning';
import './style.css';

const h=React.createElement;

const navigation=[['home','Home','House'],['voices','My Voices','Fingerprint'],['history','History','Clock3'],['settings','Settings','Settings2']];

const languages:Record<string,string>={en:'English',zh:'Chinese',ja:'Japanese',ko:'Korean',de:'German',fr:'French',es:'Spanish',pt:'Portuguese',it:'Italian',ru:'Russian',ar:'Arabic',nl:'Dutch',pl:'Polish',tr:'Turkish',vi:'Vietnamese',id:'Indonesian',th:'Thai',uk:'Ukrainian',fa:'Persian',he:'Hebrew',hu:'Hungarian',cs:'Czech',da:'Danish',sv:'Swedish',el:'Greek'};

function App(){

 const [page,setPage]=useState(location.hash.slice(1)||'home'),[fields,setFields]=useState<any>({script:'',projectName:'Untitled speech',voice:'alba',model:'pocket18',language:'en',speed:1,pitch:0,volume:1,pause:0.3,threads:2,cloneModel:'moss',cloneLanguage:'en',voiceFilter:'all',voiceLanguage:'all'});

 const [me,setMe]=useState<any>(null),[boot,setBoot]=useState<any>({}),[voices,setVoices]=useState<any[]>([]),[models,setModels]=useState<any[]>([]),[history,setHistory]=useState<any[]>([]),[settings,setSettings]=useState<any>({}),[defaultRules,setDefaultRules]=useState<any[]>([]),[driveProfile,setDriveProfile]=useState<any>(null),[driveFiles,setDriveFiles]=useState<any[]>([]),[shares,setShares]=useState<any[]>([]);

 const [error,setError]=useState(''),[modal,setModal]=useState(''),[selected,setSelected]=useState(''),[job,setJob]=useState<any>(null),[previewJob,setPreviewJob]=useState<any>(null),[segmentMode,setSegmentMode]=useState(false),[segments,setSegments]=useState<any[]>([]),[showDialogue,setDialogue]=useState(false),[showDictionary,setDictionary]=useState(false),[showBatch,setBatch]=useState(false),[transformed,setTransformed]=useState(''),[saved,setSaved]=useState('Saved'),[reference,setReference]=useState<File|null>(null),[referenceURL,setReferenceURL]=useState(''),[recording,setRecording]=useState(false),[cloneStatus,setCloneStatus]=useState(''),[cloneCreated,setCloneCreated]=useState(''),[savingClone,setSavingClone]=useState(false),[shareResult,setShareResult]=useState('');

 const [modelsLoaded,setModelsLoaded]=useState(false),[modelsError,setModelsError]=useState(''),[installSubmitting,setInstallSubmitting]=useState(false),[capturePending,setCapturePending]=useState(false);
 const jobEditorSignature=useRef('');
 const item=useRef(''),state=useRef<any>({}),project=useRef(''),undo=useRef<string[]>([]),redo=useRef<string[]>([]),fileInput=useRef<HTMLInputElement>(null),recorder=useRef<MediaRecorder|null>(null),batchIDs=useRef<string[]>([]),stream=useRef<MediaStream|null>(null);

 state.current={fields,me,job,previewJob,segments,segmentMode,showDialogue,defaultRules,models};

 const go=(p:string)=>{if(p==='dialogue'){if(!state.current.segments.length)split(true);else setDialogue(true);p='speech'};if(p==='batch'){setBatch(true);p='speech'};setPage(p);location.hash=p;setError('');document.body.classList.remove('nav-open')};

 const setField=(name:string,value:any)=>{if(name==='script'){undo.current.push(state.current.fields.script);if(undo.current.length>100)undo.current.shift();redo.current=[]};if(['script','voice','model','language'].includes(name)){setSegments([]);setSegmentMode(false)};if(['cloneModel','cloneLanguage'].includes(name))setCloneCreated('');if(name.startsWith('text_')||name.startsWith('voice_')){const next={...state.current.fields,[name]:value};next.script=editedSegments(next,state.current.segments).map(s=>s.text).join('\n');setFields(next);setSegmentMode(true);return};setFields((old:any)=>({...old,[name]:value,...(name==='search'?{historySearch:value}:{})}));if(name==='search'&&String(value).trim())go('history');if(name==='model'){setFields((old:any)=>({...old,language:'en',voice:value==='moss'?voices.find(v=>v.engine==='moss'&&!v.preset)?.id||'moss-open':'alba'}))}};

 async function api(path:string,method='GET',body?:any){const session=state.current.me;const res=await fetch('./api/'+path,{method,headers:{...(body instanceof FormData?{}:{'Content-Type':'application/json'}),'X-CSRF-Token':session?.csrf||'','X-Workspace-Expected-Account':session?.user?.id||'','X-Setup-Token':location.hash.startsWith('#setup/')?location.hash.slice(7):''},body:body===undefined?undefined:body instanceof FormData?body:JSON.stringify(body)});let data:any;try{data=await res.json()}catch{throw Error('The server could not finish this request')};if(!res.ok){if(res.status===401&&path!=='login')setMe(null);throw Error(data.error||'Request failed')};const actor=res.headers.get('X-Workspace-Account');if(session&&actor&&(actor!==session.user.id||state.current.me?.user.id!==session.user.id))throw Error('Your account changed; reload the app');return data}

 async function refreshModels(){try{const m=await api('models');setModels(m);setModelsLoaded(true);setModelsError('');return m}catch(e:any){setModelsLoaded(true);setModelsError(e.message);throw e}}
 async function ensureCloneReady(){const id=state.current.fields.cloneModel;const model=await requireCloningModel(refreshModels,id);if(state.current.fields.cloneModel!==id)throw Error('The cloning model changed. Check it is ready and try again.');return model}
 async function refresh(){const [v,m,hist,s,d]=await Promise.all([api('voices'),refreshModels(),api('history'),api('settings'),api('dictionaries')]);setVoices(v);setHistory(hist);setSettings(s);setDefaultRules(d);setFields((f:any)=>({...f,threads:s.threads}))}

 useEffect(()=>{let active=true;(async()=>{const b=await api('bootstrap');if(!active)return;setBoot(b);try{const session=await api('me');state.current.me=session;setMe(session);await refresh();const projects=await api('projects');if(projects[0]){const p=projects[0],draft=loadDraft(p);project.current=p.id;setFields((f:any)=>({...f,...draft.fields}));setSegments(draft.segments);setSegmentMode(draft.segmentMode);setDialogue(false);setSaved('Saved')}}catch{}})().catch(e=>setError(e.message));const onHash=()=>setPage(location.hash.slice(1)||'home');window.addEventListener('hashchange',onHash);return()=>{active=false;window.removeEventListener('hashchange',onHash);stream.current?.getTracks().forEach(t=>t.stop())}},[]);

 useEffect(()=>{if(!me)return;const t=setInterval(()=>{api('history?q='+encodeURIComponent(state.current.fields.historySearch||state.current.fields.search||'')).then(setHistory).catch(e=>setError(e.message));if(page==='clone'||models.some(modelInstalling))refreshModels().catch(()=>{});const current=state.current.job;if(current&&['queued','running'].includes(current.status))api('tts/jobs/'+current.id).then(j=>{if(state.current.job?.id!==j.id)return;setJob(j);if(j.status==='complete'&&scriptSignature(script())===jobEditorSignature.current){setSegments(j.script.segments);setFields((f:any)=>({...f,...segmentFields(j.script.segments)}))}}).catch(e=>setError(e.message));const preview=state.current.previewJob;if(preview&&['queued','running'].includes(preview.status))api('tts/jobs/'+preview.id).then(j=>{if(state.current.previewJob?.id===j.id)setPreviewJob(j)}).catch(e=>setError(e.message))},1600);return()=>clearInterval(t)},[me,models,page]);

 useEffect(()=>{if(!me||(!fields.script&&!project.current))return;setSaved('Saving…');const timer=setTimeout(()=>{const s=script();api('projects','POST',{...s,id:project.current}).then(p=>{project.current=p.id;setSaved('Saved')}).catch(()=>setSaved('Not saved — retry by editing'))},900);return()=>clearTimeout(timer)},[fields,segments,segmentMode,me]);

 useEffect(()=>{function key(e:KeyboardEvent){if((e.ctrlKey||e.metaKey)&&e.key==='Enter'){e.preventDefault();action('generate')}};document.addEventListener('keydown',key);return()=>document.removeEventListener('keydown',key)},[]);

 useEffect(()=>{if(!reference)return;const url=URL.createObjectURL(reference);setReferenceURL(url);return()=>URL.revokeObjectURL(url)},[reference]);

 const rules=()=>String(state.current.fields.rules||'').split('\n').filter(x=>x.includes('=')).map(x=>{const n=x.indexOf('=');return {from:x.slice(0,n).trim(),to:x.slice(n+1).trim()}}).filter(x=>x.from);

 function script(){return {...buildScript(state.current.fields,state.current.segments,state.current.segmentMode,rules()),id:project.current}}

 async function openJob(j:any){const draft=loadDraft(j.script);project.current=j.script.id||'';setJob(j);setPreviewJob(null);setSegments(draft.segments);setSegmentMode(draft.segmentMode);setDialogue(false);setFields((f:any)=>({...f,...draft.fields}));jobEditorSignature.current=scriptSignature(buildScript(draft.fields,draft.segments,draft.segmentMode,j.script.dictionary||[]));go('speech')}

 function split(customize=true){const seg=splitEditorText(state.current.fields.script,state.current.fields.voice);setSegments(seg);setFields((f:any)=>({...f,...segmentFields(seg)}));if(customize)setSegmentMode(true);setDialogue(true)}


 async function action(id:string){try{

  if(id==='navigate')id=item.current;

  if(['home','speech','clone','voices','history','settings','dialogue','batch'].includes(id)){go(id);return}

  if(id==='dismiss'){setError('');return}if(id==='menu'){document.body.classList.toggle('nav-open');return}

  if(id==='start'){setField('script',fields.idea||'');go('speech')}

  if(id==='new-project'){project.current='';setField('script','');setField('projectName','Untitled speech');setJob(null);setPreviewJob(null);setSegments([]);setDialogue(false);setSegmentMode(false)}

  if(id==='login'){await api(boot.needsSetup?'setup':'login','POST',{email:fields.email,password:fields.password});const session=await api('me');state.current.me=session;setMe(session);setField('password','');if(boot.needsSetup){setBoot({...boot,needsSetup:false});go('home')};await refresh()}

  if(id==='logout'){const res=await fetch('./auth/logout',{method:'POST',headers:{'X-CSRF-Token':me.csrf}});if(!res.ok)throw Error('Sign out failed');setModal('');setMe(null);setVoices([]);setHistory([]);setJob(null)}

  if(id==='generate'){const current=script();jobEditorSignature.current=scriptSignature(current);const payload={...current,requestId:crypto.randomUUID()};let j;try{j=await api('tts/jobs','POST',payload)}catch(e){if(!(e instanceof TypeError))throw e;j=await api('tts/jobs','POST',payload)};setJob(j);go('speech')}

  if(id==='cancel'){setJob(await api('tts/jobs/'+job.id+'/cancel','POST',{}))}

  if(id==='retry'){setJob(await api('tts/jobs/'+job.id+'/retry','POST',{}))}

  if(id==='undo'&&undo.current.length){redo.current.push(fields.script);setSegments([]);setSegmentMode(false);setFields((f:any)=>({...f,script:undo.current.pop()}))}

  if(id==='redo'&&redo.current.length){undo.current.push(fields.script);setSegments([]);setSegmentMode(false);setFields((f:any)=>({...f,script:redo.current.pop()}))}

  if(id==='import')fileInput.current?.click();

  if(id==='dictionary')setDictionary(!showDictionary);

  if(id==='toggle-dialogue'){if(!showDialogue&&!segments.length)split(false);else setDialogue(!showDialogue)}

  if(id==='split')split();

  if(id==='preview-text'){let t=fields.script;const effective=new Map([...defaultRules,...rules()].map(r=>[r.from,r.to]));for(const [from,to] of effective){const parts=t.split(from);if(t.length+(parts.length-1)*(to.length-from.length)>400000)throw Error('Pronunciation preview is too long; reduce the replacements');t=parts.join(to)}setTransformed(t)}

  if(id==='save-dictionary'){await api('dictionaries','POST',rules());setDefaultRules(rules());setSaved('Dictionary saved')}

  if(id==='regenerate'){const s=buildScript(fields,segments,true,rules());setSegmentMode(true);jobEditorSignature.current=scriptSignature(s);setJob(await api('tts/jobs','POST',{...s,parent:job.id,replace:item.current}))}

  if(id==='preview-voice'||id==='listen-voice'){const voice=id==='listen-voice'?item.current:fields.voice;const j=await api('voices/'+voice+'/preview','POST',{});setPreviewJob(j)}

  if(id==='choose-voice'){setSegments([]);setSegmentMode(false);const v=voices.find(v=>v.id===item.current);setFields((f:any)=>({...f,voice:v.id,model:v.engine==='moss'?'moss':v.engine==='pocket'?'pocket18':v.engine,language:v.language}));go('speech')}

  if(id==='open-job'){const j=history.find(j=>j.id===item.current);if(j)await openJob(j)}

  if(id==='save-settings'){await api('settings','PATCH',{threads:Number(fields.threads)});await refresh();setSaved('Settings saved')}

  if(id==='refresh-models'){await refreshModels();setError('');return}

  if(id==='install-cloning-model'||id==='install-model'){if(me?.user?.role!=='owner')throw Error('Ask the Workspace owner to install the cloning model on this screen.');const idToInstall=id==='install-cloning-model'?state.current.fields.cloneModel:item.current;if(!state.current.models.some((m:any)=>m.id===idToInstall))throw Error('Refresh the model list before installing');setSelected(idToInstall);setField('hfToken','');setField('license',false);setField('enableCloning',false);setModal('install')}

  if(id==='confirm-install'){if(installSubmitting)return;setInstallSubmitting(true);try{await api('models/'+selected+'/install','POST',{license:!!fields.license,token:fields.hfToken||'',cloning:!!fields.enableCloning});setField('hfToken','');setModal('');await refreshModels()}finally{setInstallSubmitting(false)}}

  if(id==='remove-model'){setSelected(item.current);setModal('remove-model')}

  if(id==='confirm-remove-model'){await api('models/'+selected,'DELETE');setModal('');await refresh()}

  if(id==='record'){

   if(recording){recorder.current?.stop();return}

   if(capturePending)return;setCapturePending(true);
   try{await ensureCloneReady();
   stream.current=await navigator.mediaDevices.getUserMedia({audio:true});const chunks:BlobPart[]=[];const r=new MediaRecorder(stream.current);recorder.current=r;r.ondataavailable=e=>chunks.push(e.data);r.onstop=()=>{stream.current?.getTracks().forEach(t=>t.stop());setReference(new File(chunks,'recording.webm',{type:r.mimeType}));setRecording(false)};r.start();setRecording(true);setTimeout(()=>{if(r.state==='recording')r.stop()},30000)
   }catch(e){stream.current?.getTracks().forEach(t=>t.stop());throw e}finally{setCapturePending(false)}

  }

  if(id==='save-clone'){if(savingClone)return;setError('');setSavingClone(true);setCloneStatus('Checking cloning model…');try{await ensureCloneReady();if(!reference)throw Error('Record or choose a reference first');const form=new FormData();form.set('audio',reference);form.set('name',fields.cloneName||'');form.set('tags',fields.cloneTags||'');form.set('consent',String(!!fields.consent));form.set('model',state.current.fields.cloneModel);form.set('language',fields.cloneLanguage||'en');setCloneStatus('Analysing reference and preparing voice…');const result=await api('voices','POST',form);setCloneCreated(result.voice.id);setCloneStatus(result.audio.duration.toFixed(1)+'s reference · '+(result.preparation?.summary||'saved privately'));await refresh()}finally{setSavingClone(false)}}

  if(id==='test-clone'){const j=await api('voices/'+cloneCreated+'/preview','POST',{});setPreviewJob(j)}
  if(id==='close-preview')setPreviewJob(null);
  if(id==='cancel-preview')setPreviewJob(await api('tts/jobs/'+previewJob.id+'/cancel','POST',{}));
  if(id==='retry-preview')setPreviewJob(await api('tts/jobs/'+previewJob.id+'/retry','POST',{}));

  if(id==='use-clone'){setSegments([]);setSegmentMode(false);const v=voices.find(v=>v.id===cloneCreated);setFields((f:any)=>({...f,voice:v.id,model:v.engine,language:v.language}));go('speech')}

  if(id==='run-batch'){const lines=String(fields.batchText||'').split('\n').map(x=>x.trim()).filter(Boolean);if(!lines.length||lines.length>20)throw Error('Add 1–20 lines');batchIDs.current=[];for(let i=0;i<lines.length;i++){const j=await api('tts/jobs','POST',{...script(),text:lines[i],segments:[],name:'Batch '+(i+1)});batchIDs.current.push(j.id)};await refresh();go('history')}

  if(id==='export-batch'||id==='export-all'){const ids=(id==='export-all'?history.filter(j=>j.status==='complete').slice(0,20).map(j=>j.id):batchIDs.current.filter(id=>history.some(j=>j.id===id&&j.status==='complete')));if(!ids.length)throw Error('Wait for a recording to complete');location.href='./api/exports/zip?ids='+ids.join(',')}

  if(id==='edit-voice'){const v=voices.find(v=>v.id===item.current);setSelected(v.id);setFields((f:any)=>({...f,editName:v.name,editTags:v.tags,favorite:v.favorite,deleteConsent:false}));setShareResult('');setShares(await api('voices/'+v.id+'/shares'));setModal('manage')}

  if(id==='save-voice'){await api('voices/'+selected,'PATCH',{name:fields.editName,tags:fields.editTags,favorite:!!fields.favorite});setModal('');await refresh()}

  if(id==='favorite-voice'){const v=voices.find(v=>v.id===item.current);await api('voices/'+v.id+'/favorite','POST',{favorite:!v.favorite});await refresh()}

  if(id==='delete-recording'){setSelected(item.current);setField('deleteRecordingConsent',false);setModal('delete-recording')}

  if(id==='confirm-delete-recording'){if(!fields.deleteRecordingConsent)throw Error('Confirm recording deletion first');await api('history/'+selected,'DELETE');setModal('');if(job?.id===selected)setJob(null);await refresh()}

  if(id==='delete-voice'){if(!fields.deleteConsent)throw Error('Confirm voice deletion first');await api('voices/'+selected,'DELETE');setModal('');await refresh()}

  if(id==='share-voice'){const grant=await api('voices/'+selected+'/share','POST',{user:fields.shareEmail});setShares(await api('voices/'+selected+'/shares'));setShareResult('Private use granted. You can revoke it below.')}

  if(id==='revoke-specific'){const result=await api('voices/'+selected+'/share/'+item.current,'DELETE');setShares(await api('voices/'+selected+'/shares'));setShareResult(result.revoked?'Access revoked.':'No matching grant.')}

  if(id==='help')setModal('help');if(id==='account')setModal('account');if(id==='close-modal')setModal('');

  if(id==='drive'){await saveDrive()}

  if(id==='import-drive'){const profile=await verifyImportDrive();setDriveProfile(profile);const scope=profile.scopes?.find((s:any)=>s.name==='My Drive')||profile.scopes?.[0];if(!scope)throw Error('Your account has no accessible Drive sources');setFields((f:any)=>({...f,driveSource:scope.name,drivePath:scope.scope||'/'}));setDriveFiles([]);setModal('drive-import');await listImportDrive(scope.name,scope.scope||'/')}

  if(id==='list-drive')await listImportDrive(fields.driveSource,fields.drivePath||'/');

  if(id==='select-drive-file'){const file=driveFiles.find(f=>f.path===item.current);if(!file)throw Error('Choose a current file');if(file.type==='directory'){setField('drivePath',file.path);await listImportDrive(fields.driveSource,file.path)}else{await verifyImportDrive();const res=await importDriveRequest('resources/download?'+new URLSearchParams({source:fields.driveSource,file:file.path}));if(Number(res.headers.get('Content-Length')||0)>200000){await res.body?.cancel();throw Error('Import a TXT or Markdown file smaller than 200 KB')};const reader=res.body!.getReader(),chunks:Uint8Array[]=[];let size=0;while(true){const part=await reader.read();if(part.done)break;size+=part.value.length;if(size>200000){await reader.cancel();throw Error('Text file exceeds 200 KB')};chunks.push(part.value)};await verifyImportDrive();const text=await new Blob(chunks as BlobPart[]).text();project.current='';setField('script',text);setField('projectName',file.name.replace(/\.(txt|md)$/i,''));setSegments([]);setDialogue(false);setSegmentMode(false);setModal('');go('speech')}}

 }catch(e:any){setError(e.message);setCloneStatus('')}}

 async function importDriveRequest(path:string){if(!boot.workspace||!boot.oidc)throw Error('My Drive import is available through One Workspace');const user=state.current.me?.user.id;const r=await fetch('/drive/api/'+path,{credentials:'same-origin'});if(state.current.me?.user.id!==user)throw Error('Your account changed; reload before importing');if(!r.ok)throw Error(r.status===401?'Open Drive from Apps with the same OneAccount, then retry.':'Drive denied this file or folder ('+r.status+')');return r}

 async function verifyImportDrive(){const profile=await(await importDriveRequest('users?username=self')).json();if(profile.username!==state.current.me?.user.subject||String(profile.loginMethod).toLowerCase()!=='oidc')throw Error('Drive must be signed in to the same OneAccount');return profile}

 async function listImportDrive(source:string,folder:string){await verifyImportDrive();const listing=await(await importDriveRequest('resources?'+new URLSearchParams({source,path:folder}))).json();if(listing.type!=='directory')throw Error('Choose a folder first');const items=[...(listing.folders||[]),...(listing.files||[])].filter(f=>f.type==='directory'||/\.(txt|md)$/i.test(f.name));if(items.some(f=>typeof f.name!=='string'||/[\/\\\u0000]/.test(f.name)||f.name==='.'||f.name==='..'))throw Error('Drive returned an invalid file name');setDriveFiles(items.map(f=>({...f,path:folder.replace(/\/$/,'')+'/'+f.name,icon:f.type==='directory'?'FolderOpen':'FileText'})))}

 async function saveDrive(){

  if(!boot.workspace||!boot.oidc)throw Error('My Drive is available when My Voice is opened through One Workspace. Download a copy in standalone mode.');

  const expected=state.current.me.user.id,subject=state.current.me.user.subject;

  async function drive(path:string,options:any={}){if(state.current.me?.user.id!==expected)throw Error('Your account changed; reload before saving');const r=await fetch('/drive/api/'+path,{credentials:'same-origin',...options});if(!r.ok)throw Object.assign(Error(r.status===401?'Open Drive from Apps to connect the same account, then retry.':'Drive returned '+r.status+'; your original recording is safely saved.'),{status:r.status});return r}

  async function verify(){const profile=await(await drive('users?username=self')).json();if(profile.username!==subject||String(profile.loginMethod).toLowerCase()!=='oidc')throw Error('Drive must be signed in to the same OneAccount');return profile}

  const profile=await verify(),scope=profile.scopes?.find((s:any)=>s.name==='My Drive')||profile.scopes?.[0];if(!scope)throw Error('Your Drive account has no writable source');

  const folder=String(scope.scope||'/').replace(/\/$/,'')+'/My Voice',path=folder+'/speech-'+job.id+'.wav',q=(path:string)=>new URLSearchParams({source:scope.name,path}).toString();

  try{await drive('resources?'+q(folder))}catch(e:any){if(e.status!==404)throw e;await verify();try{await drive('resources?'+q(folder)+'&isDir=true&override=false',{method:'POST'})}catch(e:any){if(e.status!==409)throw e}}

  const audio=await fetch('./api/outputs/'+job.id,{headers:{'X-Workspace-Expected-Account':expected}});if(!audio.ok)throw Error('Recording is no longer available');const blob=await audio.blob();await verify();

  try{await drive('resources?'+q(path)+'&override=false',{method:'POST',headers:{'Content-Type':'audio/wav'},body:blob})}catch(e:any){if(e.status!==409)throw e}

  // Reconcile a lost response or retry by content hash; never overwrite a file.

  const remote=await(await drive('resources/download?'+new URLSearchParams({source:scope.name,file:path}))).arrayBuffer();const local=await blob.arrayBuffer();const digest=async(b:ArrayBuffer)=>Array.from(new Uint8Array(await crypto.subtle.digest('SHA-256',b))).join(',');if(await digest(remote)!==await digest(local))throw Error('Drive contains a different file at this destination; no existing file was replaced');await verify();setError('Saved to My Drive · '+scope.name+path);

 }

 const selectedVoice=voices.find(v=>v.id===fields.voice),model=models.find(m=>m.id===fields.model),install=models.find(m=>m.id===selected);

 const voiceOptions=voices.filter(v=>fields.model==='moss'?v.engine==='moss':v.engine==='pocket'||v.engine===fields.model).map(v=>({value:v.id,label:v.name+(v.preset?'':' · Clone')}));

 const cloneSetup=cloningSetup(models,fields.cloneModel,modelsLoaded,modelsError,me?.user?.role==='owner');
 useEffect(()=>{if(cloneSetup.model&&cloneSetup.model.id!==fields.cloneModel)setField('cloneModel',cloneSetup.model.id)},[models,fields.cloneModel]);
 const data:any={loggedOut:!me,signInLabel:boot.needsSetup?'Create your local account':'Sign in',oidc:boot.oidc,error,workspace:boot.workspace||'/',version:'My Voice '+(boot.version||''),initials:me?.user?.name?.slice(0,1)?.toUpperCase()||'M',nav:navigation.map(([id,label,icon])=>({id,label,icon,variant:page===id?'nav-active':'nav'})),pageTitle:({home:'Home',speech:'Text to Speech',clone:'Voice Cloning',voices:'My Voices',history:'History',settings:'Settings'} as any)[page]||'Home',pageIcon:({home:'House',speech:'AudioLines',clone:'Fingerprint',voices:'Library',history:'Clock3',settings:'Settings2'} as any)[page],home:page==='home',speech:page==='speech',clone:page==='clone',voicesPage:page==='voices',historyPage:page==='history',settingsPage:page==='settings',quickTools:[{id:'speech',icon:'AudioLines',label:'Speech'},{id:'clone',icon:'Fingerprint',label:'Voice Clone'},{id:'dialogue',icon:'MessagesSquare',label:'Dialogue'},{id:'batch',icon:'Layers',label:'Batch'},{id:'history',icon:'FolderOpen',label:'My Creations'}],starts:[{id:'speech',icon:'Feather',label:'Tell your story',detail:'Turn your script into natural speech'},{id:'clone',icon:'Fingerprint',label:'Find your sound',detail:'Create a voice that feels like you'},{id:'dialogue',icon:'MessagesSquare',label:'Start a conversation',detail:'Bring different characters to life'}],hasHistory:history.length>0,recent:history.slice(0,4).map(j=>({...j,detail:j.status+' · '+new Date(j.created).toLocaleDateString()})),filteredHistory:history.filter(j=>(j.script.name+' '+j.script.text).toLowerCase().includes(String(fields.historySearch||fields.search||'').toLowerCase())).map(j=>({...j,detail:new Date(j.created).toLocaleString()+' · '+(j.duration?j.duration.toFixed(1)+'s':'Queued')})),saved,counter:fields.script.length.toLocaleString()+' characters · '+Math.ceil(fields.script.split(/\s+/).filter(Boolean).length/150*60)+'s estimated',showDialogue,showDictionary,showBatch,segmentDetailsLabel:showDialogue?'Hide segment details':'Show segment details',dialogueVariant:showDialogue?'selected':'text',hasPreview:!!previewJob&&!(page==='clone'&&previewJob.script.voice===cloneCreated),previewName:previewJob?.script?.name,previewStatus:previewJob?.error||previewJob?.status,previewRunning:['queued','running'].includes(previewJob?.status),previewFailed:['failed','cancelled'].includes(previewJob?.status),voiceOptions,voiceName:selectedVoice?.name||'Choose a voice',voiceDescription:selectedVoice?.tags||'',modelOptions:[{value:'pocket18',label:'Pocket TTS · 18.2M · Fast'},{value:'pocket7',label:'Pocket TTS · 6.9M · Lightest'},{value:'moss',label:'MOSS TTS Nano · Cloning & multilingual'}],modelDescription:model?.installed?'Installed · CPU only · '+model.version:'Install this model in Settings to generate speech.',languages:(model?.languages||['en']).map((v:string)=>({value:v,label:languages[v]||v})),speedLabel:'Speed · '+Number(fields.speed).toFixed(2)+'×',pitchLabel:'Pitch · '+fields.pitch+' semitones',volumeLabel:'Volume · '+Math.round(fields.volume*100)+'%',pauseLabel:'Segment pause · '+Number(fields.pause).toFixed(1)+'s',generating:job&&['queued','running'].includes(job.status),defaultRuleText:defaultRules.length?'Account defaults: '+defaultRules.map(r=>r.from+' → '+r.to).join('; '):'No account defaults saved.',segments:segments.map((s,i)=>({...s,number:String(i+1).padStart(2,'0'),label:'Segment '+(i+1),textField:'text_'+s.id,voiceField:'voice_'+s.id,disabled:job?.status!=='complete'||segments.length!==job?.script?.segments?.length||segments.some((s,i)=>s.id!==job?.script?.segments?.[i]?.id)})),transformed,hasJob:!!job,jobName:job?.script?.name,jobStatus:job?.error||((job?.status||'')+' · '+(job?.status==='complete'?job.duration.toFixed(1)+'s · '+job.wall.toFixed(1)+'s generation':(job?.progress||0)+'%')),jobRunning:job&&['queued','running'].includes(job.status),jobComplete:job?.status==='complete',jobFailed:['cancelled','failed'].includes(job?.status),wav:'./api/outputs/'+job?.id+'/download?format=wav',mp3:'./api/outputs/'+job?.id+'/download?format=mp3',recordLabel:recording?'Stop recording':capturePending?'Preparing to record…':'Record your voice',cloneOptions:cloneSetup.options,cloneReady:cloneSetup.ready,cloneSetupTitle:cloneSetup.title,cloneSetupDetail:cloneSetup.detail,cloneSetupStatus:cloneSetup.status,cloneCanInstall:cloneSetup.canInstall,cloneInstallLabel:cloneSetup.installLabel,cloneNeedsOwner:cloneSetup.needsOwner,cloneSetupBusy:cloneSetup.busy,cloneRecordDisabled:!recording&&(!cloneSetup.ready||capturePending||savingClone),cloneCreateDisabled:!cloneSetup.ready||savingClone||!reference||!fields.cloneName?.trim()||!fields.consent,cloneModelDisabled:recording||capturePending||savingClone||cloneSetup.busy,installDisabled:installSubmitting||!fields.license,cloneLanguages:(models.find(m=>m.id===fields.cloneModel)?.languages||['en']).map((l:string)=>({value:l,label:languages[l]||l})),voiceLanguages:[{value:'all',label:'All reference languages'},...Array.from(new Set(voices.map(v=>v.language))).map((l:any)=>({value:l,label:languages[l]||l}))],cloneStatus,cloneCreated:!!cloneCreated,savingClone,voiceFilters:[{value:'all',label:'All voices'},{value:'clones',label:'My clones'},{value:'favourites',label:'Favourites'}],filteredVoices:voices.filter(v=>(v.name+' '+v.tags).toLowerCase().includes(String(fields.voiceSearch||'').toLowerCase())&&(fields.voiceFilter!=='clones'||!v.preset)&&(fields.voiceFilter!=='favourites'||v.favorite)&&(fields.voiceLanguage==='all'||v.language===fields.voiceLanguage)).map(v=>({...v,type:v.preset?'PRESET':v.shared?'SHARED':'MY VOICE',favoriteLabel:v.favorite?'Remove favourite':'Favourite voice',favoriteVariant:v.favorite?'icon favourite':'icon',owned:!v.preset&&!v.shared,engineLabel:v.engine==='moss'?'MOSS · CPU':v.engine==='pocket'?'Pocket · English':v.engine+' · model-specific'})),models:models.map(m=>({...m,installLabel:m.installed?'Installed':'Not installed',buttonLabel:m.installed?'Reinstall / change access':'Download model',description:(m.size/1024/1024).toFixed(0)+' MB download · '+m.languages.map((l:string)=>languages[l]||l).join(', ')})),storageInfo:(settings.storage||'')+' · '+((settings.usedBytes||0)/1073741824).toFixed(2)+' GB used · '+((settings.availableBytes||0)/1073741824).toFixed(1)+' GB available',retentionInfo:settings.retention||'',modal:!!modal,modalTitle:({install:'Install '+install?.name,manage:'Manage your voice',help:'A little help, a lot of possibilities',account:'Your account','remove-model':'Remove model','delete-recording':'Remove recording','drive-import':'Import from My Drive'} as any)[modal],installModal:modal==='install',manageModal:modal==='manage',helpModal:modal==='help',accountModal:modal==='account',pocketInstall:false,driveImportModal:modal==='drive-import',driveFiles,driveSources:(driveProfile?.scopes||[]).map((s:any)=>({value:s.name,label:s.name})),deleteRecordingModal:modal==='delete-recording',removeModelModal:modal==='remove-model',installDescription:install?.runtime+' · '+Math.ceil((install?.size||0)/1024/1024)+' MB. Downloads are checksum verified. No model weights are bundled with the launcher.',installURL:install?.url,shares,shareResult,accountName:me?.user?.name||me?.user?.email};

 const widgets:any={

 'preview-audio':previewJob?.status==='complete'?h('audio',{key:previewJob.id,controls:true,src:'./api/outputs/'+previewJob.id,preload:'metadata','aria-label':'Voice preview audio'}):null,

 'feature-wave':h('div',{className:'feature-wave','aria-hidden':true},Array.from({length:52},(_,i)=>h('i',{key:i,style:{height:(20+Math.sin(i*.55)**2*70+Math.sin(i*.12)**2*90)+'px'}}))),

 waveform:job?.waveform?.length?h('svg',{viewBox:'0 0 640 80',className:'waveform',role:'img','aria-label':'Generated audio waveform'},job.waveform.map((v:number,i:number)=>h('rect',{key:i,x:i*4,y:40-Math.max(2,v*38),width:2.6,height:Math.max(4,v*76),rx:1.3}))):job?h('progress',{max:100,value:job.progress||0,'aria-label':'Generation progress'}):null,

 audio:job?.status==='complete'?h('audio',{key:job.id,controls:true,src:'./api/outputs/'+job.id,preload:'metadata','aria-label':'Generated speech'}):null,

 upload:h('label',{className:'upload-area'+(!cloneSetup.ready?' upload-area--disabled':''),'aria-disabled':!cloneSetup.ready},h('input',{type:'file',disabled:!cloneSetup.ready||recording||capturePending||savingClone,accept:'.wav,.mp3,.m4a,.webm,.ogg',onChange:async(e:any)=>{const file=e.target.files?.[0];if(!file)return;try{await ensureCloneReady();if(file.size>25*1024*1024)throw Error('Choose a reference smaller than 25 MB');setReference(file);setCloneCreated('');setError('')}catch(err:any){e.target.value='';setError(err.message)}}}),h('span',{className:'upload-icon'},'↑'),h('strong',{},reference?.name||(cloneSetup.ready?'Drop in your voice':'Install the cloning model first')),h('span',{},cloneSetup.ready?'or click to choose a reference recording':'Recording and upload unlock when the model is ready')),

 'clone-model-progress':cloneSetup.busy?h('progress',{max:100,value:cloneSetup.progress,'aria-label':cloneSetup.progress===undefined?'Preparing cloning model':'Cloning model download progress'}):null,

 'reference-audio':referenceURL?h('audio',{controls:true,src:referenceURL,'aria-label':'Original voice reference'}):null,

 'prepared-reference-audio':cloneCreated?h('audio',{key:cloneCreated,controls:true,src:'./api/voices/'+cloneCreated+'/reference','aria-label':'Prepared voice reference'}):null,

 'clone-test-audio':previewJob?.script?.voice===cloneCreated?(previewJob.status==='complete'?h('audio',{controls:true,src:'./api/outputs/'+previewJob.id,'aria-label':'Cloned voice test sample'}):h('span',{className:'text text--caption'},previewJob.error||'Preparing test sample · '+previewJob.status)):null

 };

 return h(Bindings.Provider,{value:{fields,setField,item,widgets,action}},h(React.Fragment,{},h(Voice,{data,onAction:action}),h('input',{type:'file',accept:'.txt,.md',hidden:true,ref:fileInput,onChange:async(e:any)=>{const file=e.target.files?.[0];if(!file)return;if(file.size>200000){setError('Import a text file smaller than 200 KB');return};setField('script',await file.text());setField('projectName',file.name.replace(/\.(txt|md)$/i,''));setSegments([]);setDialogue(false);setSegmentMode(false)}})));

}

createRoot(document.getElementById('root')!).render(h(App));
