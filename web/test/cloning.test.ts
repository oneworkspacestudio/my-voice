import {test} from 'node:test';
import assert from 'node:assert/strict';
import {cloningSetup,requireCloningModel,cloneReady} from '../src/cloning.ts';

const moss={id:'moss',name:'MOSS TTS Nano',cloning:true,installed:false,cloningReady:false,size:763682544,languages:['en'],status:''};
const pocket={...moss,id:'pocket18',name:'Pocket',cloning:false,installed:true};

test('an uninstalled cloning model remains discoverable and installable before recording',()=>{
 const s=cloningSetup([pocket,moss],'moss',true,'',true);
 assert.deepEqual(s.options,[{value:'moss',label:'MOSS TTS Nano'}]);
 assert.equal(s.canInstall,true);assert.equal(s.ready,false);assert.match(s.title,/before recording/);
});

test('downloading and reference preparation stay locked until installation is verified',()=>{
 for(const status of ['Downloading 100%','Preparing cloning model…']){
  const s=cloningSetup([{...moss,status}],'moss',true,'',true);
  assert.equal(s.busy,true);assert.equal(s.canInstall,false);assert.equal(s.ready,false);
 }
 const s=cloningSetup([{...moss,installed:true,cloningReady:true,status:'Ready'}],'moss',true,'',true);
 assert.equal(s.ready,true);assert.equal(s.canInstall,false);
});

test('setup failures can be retried and a failed refresh never unlocks capture',()=>{
 const s=cloningSetup([{...moss,status:'Download failed. Check your connection and retry'}],'moss',true,'',true);
 assert.equal(s.installLabel,'Retry model setup');assert.equal(s.canInstall,true);
 assert.equal(cloningSetup([{...moss,installed:true,cloningReady:true}],'moss',true,'offline',true).ready,false);
});

test('loading, empty catalog and non-owner states explain the next step',()=>{
 assert.match(cloningSetup([],'moss',false,'',true).title,/Checking/);
 assert.match(cloningSetup([],'moss',true,'',true).title,/No cloning model/);
 const s=cloningSetup([moss],'moss',true,'',false);assert.equal(s.canInstall,false);assert.equal(s.needsOwner,true);
});

test('fresh readiness check rejects a removed model and permits a verified model',async()=>{
 await assert.rejects(requireCloningModel(async()=>[moss],'moss'),/wait until it is ready/);
 await assert.rejects(requireCloningModel(async()=>[],'moss'),/Install/);
 const ready={...moss,installed:true,cloningReady:true};assert.equal(await requireCloningModel(async()=>[ready],'moss'),ready);
 assert.equal(cloneReady({...ready,status:'Preparing cloning model…'}),false);
});
