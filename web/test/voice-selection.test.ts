import {test} from 'node:test';
import assert from 'node:assert/strict';
import {speechVoiceOptions,selectSpeechVoice,selectSpeechModel,recordingVoiceLabel,recordingTitle,voiceFitsModel} from '../src/voice-selection.ts';
import {buildScript,loadDraft} from '../src/editor.ts';

const voices=[
 {id:'alba',name:'Alba',engine:'pocket',language:'en',preset:true},
 {id:'moss-open',name:'Open Voice',engine:'moss',language:'en',preset:true},
 {id:'own',name:'My own voice',engine:'moss',language:'en'},
 {id:'shared',name:'Shared speaker',engine:'moss',language:'fr',shared:true},
 {id:'pocket-clone',name:'Pocket clone',engine:'pocket18',language:'en'}
];
const models=[{id:'pocket18',languages:['en']},{id:'pocket7',languages:['en']},{id:'moss',languages:['en','fr','ja']}];
const original={name:'Voice preview · Alba',text:'Welcome. This ending must be spoken too.',voice:'alba',model:'pocket18',language:'en',inputMode:'text',segments:[{id:'cached',text:'Welcome.',voice:'alba'}]};

test('a clone is directly selectable from an Alba preview and the complete payload uses its model and name',()=>{
 const draft=loadDraft(original);
 assert.ok(speechVoiceOptions(voices).some(v=>v.value==='own'));
 assert.equal(speechVoiceOptions(voices)[0].value,'own');
 const selected=selectSpeechVoice(draft.fields,'own',voices,models);
 const payload=buildScript(selected,[],false,[]);
 assert.equal(payload.voice,'own');assert.equal(payload.model,'moss');
 assert.equal(payload.name,'Voice preview · My own voice');assert.equal(payload.text,original.text);
 assert.equal(payload.inputMode,'text');assert.deepEqual(payload.segments,[]);
 assert.equal(recordingVoiceLabel(payload,voices),'Voice · My own voice');
 assert.equal(original.name,'Voice preview · Alba');
 assert.equal(recordingVoiceLabel(original,voices),'Voice · Alba');
});

test('voice selection preserves a custom title and a language supported by the new model',()=>{
 const fields={...loadDraft(original).fields,projectName:'A title I wrote',language:'fr'};
 const selected=selectSpeechVoice(fields,'own',voices,models);
 assert.equal(selected.projectName,fields.projectName);assert.equal(selected.language,'fr');
 assert.equal(selected.script,fields.script);
});

test('an existing clone draft with a stale automatic Alba title is corrected without changing its recording',()=>{
 const saved={...original,voice:'own',model:'moss'};
 const selected=selectSpeechVoice(loadDraft(saved).fields,'own',voices,models);
 assert.equal(selected.projectName,'Voice preview · My own voice');
 assert.equal(recordingTitle(saved,voices),'Voice preview · My own voice');
 assert.equal(saved.name,'Voice preview · Alba');
 assert.equal(recordingTitle({...saved,name:'My chosen title'},voices),'My chosen title');
});

test('switching between clone and preset keeps the model and language compatible',()=>{
 let fields=selectSpeechVoice({...loadDraft(original).fields,language:'ja'},'own',voices,models);
 fields=selectSpeechVoice(fields,'alba',voices,models);
 assert.equal(fields.voice,'alba');assert.equal(fields.model,'pocket18');assert.equal(fields.language,'en');
 assert.equal(fields.projectName,'Voice preview · Alba');
 assert.equal(selectSpeechVoice({...fields,model:'pocket7'},'alba',voices,models).model,'pocket7');
 assert.equal(selectSpeechVoice(fields,'pocket-clone',voices,models).model,'pocket18');
});

test('reselecting a compatible model keeps an explicitly selected clone; incompatible model uses its preset',()=>{
 const clone=selectSpeechVoice(loadDraft(original).fields,'own',voices,models);
 assert.equal(selectSpeechModel(clone,'moss',voices,models).voice,'own');
 const pocket=selectSpeechModel(clone,'pocket7',voices,models);
 assert.equal(pocket.voice,'alba');assert.equal(pocket.model,'pocket7');
 const moss=selectSpeechModel(pocket,'moss',voices,models);
 assert.equal(moss.voice,'moss-open');assert.equal(moss.model,'moss');
});

test('shared voices are included and stale deleted selections never silently substitute another speaker',()=>{
 assert.ok(speechVoiceOptions(voices).some(v=>v.value==='shared'&&v.label.includes('Shared')));
 const fields=loadDraft(original).fields;
 assert.equal(selectSpeechVoice(fields,'deleted',voices,models),fields);
 assert.equal(selectSpeechVoice(fields,'shared',voices,models).model,'moss');
 assert.ok(!voiceFitsModel(voices[2],'pocket18'));
 assert.ok(voiceFitsModel(voices[2],'moss'));
});

test('recording labels describe the immutable recorded speakers, including explicit dialogue',()=>{
 assert.equal(recordingVoiceLabel({...original,inputMode:'segments',segments:[{voice:'own'},{voice:'shared'},{voice:'own'}]},voices),'Voices · My own voice, Shared speaker');
 assert.equal(recordingVoiceLabel({...original,voice:'removed'},voices),'Voice · Unavailable voice');
});

test('same-name clone variants show their actual model and select its matching language',()=>{
 const variants=[...voices,{id:'q06',name:'My own voice',engine:'qwen06',language:'en'},{id:'q17',name:'My own voice',engine:'qwen17',language:'en'}];
 const choices=[...models,{id:'qwen06',shortName:'Qwen3 0.6B',languages:['en','zh']},{id:'qwen17',shortName:'Qwen3 1.7B',languages:['en','zh']}];
 const options=speechVoiceOptions(variants,choices);
 assert.match(options.find(v=>v.value==='q06')!.label,/Qwen3 0.6B/);
 assert.match(options.find(v=>v.value==='q17')!.label,/Qwen3 1.7B/);
 let selected=selectSpeechVoice({...loadDraft(original).fields,language:'zh'},'q17',variants,choices);
 assert.equal(selected.voice,'q17');assert.equal(selected.model,'qwen17');assert.equal(selected.language,'zh');
 assert.match(recordingVoiceLabel(buildScript(selected,[],false,[]),variants,choices),/Qwen3 1.7B$/);
 selected=selectSpeechVoice(selected,'own',variants,choices);
 assert.equal(selected.model,'moss');assert.equal(selected.language,'en');
});
