import {test} from 'node:test';
import assert from 'node:assert/strict';
import {buildScript, loadDraft, splitEditorText, scriptSignature} from '../src/editor.ts';

const fields = {projectName:'Story',script:'First sentence. Second sentence. Final sentence.',voice:'alba',model:'pocket18',language:'en',speed:1,pitch:0,volume:1,pause:.3};
const cached = [{id:'old',text:'First sentence.',voice:'alba'}];

test('the entire edited script overrides cached chunks from an opened recording',()=>{
  const draft=loadDraft({...fields,name:'Story',text:'First sentence.',segments:cached});
  const payload=buildScript({...draft.fields,script:fields.script},draft.segments,draft.segmentMode,[]);
  assert.equal(payload.inputMode,'text');
  assert.equal(payload.text,fields.script);
  assert.deepEqual(payload.segments,[]);
});

test('viewing segment details does not select a different generation source',()=>{
  const segments=splitEditorText(fields.script,'alba');
  assert.deepEqual(buildScript(fields,segments,false,[]),buildScript(fields,[],false,[]));
});

test('a legacy saved project keeps its full editor text and discards mismatched cached details',()=>{
  const draft=loadDraft({...fields,name:'Story',text:fields.script,segments:cached});
  assert.equal(draft.fields.script,fields.script);
  assert.deepEqual(draft.segments,[]);
  assert.equal(draft.segmentMode,false);
});

test('edited multi-voice segments retain every segment and canonical full text',()=>{
  const segments=splitEditorText(fields.script,'alba');
  const edited={...fields,['text_'+segments[2].id]:'The ending has changed.',['voice_'+segments[1].id]:'another-voice'};
  const payload=buildScript(edited,segments,true,[]);
  assert.equal(payload.inputMode,'segments');
  assert.equal(payload.segments.length,3);
  assert.equal(payload.segments[1].voice,'another-voice');
  assert.equal(payload.text,'First sentence.\nSecond sentence.\nThe ending has changed.');
  assert.equal(loadDraft(payload).segmentMode,true);
});

test('long Unicode and punctuation survive chunking without exceeding native limits',()=>{
  const text='... '+('Café story 中文 words! '.repeat(100))+' THE FINAL ENDING?!';
  const segments=splitEditorText(text,'alba');
  assert.ok(segments.length>5);
  assert.equal(segments.map(s=>s.text).join('').replace(/\s/gu,''),text.replace(/\s/gu,''));
  assert.ok(segments.every(s=>new TextEncoder().encode(s.text).length<=500));
});

test('asynchronous metadata does not invalidate an unchanged editor; text changes do',()=>{
  const script=buildScript(fields,[],false,[]);
  assert.equal(scriptSignature(script),scriptSignature({...script,id:'saved-project'} as typeof script));
  assert.notEqual(scriptSignature(script),scriptSignature({...script,text:script.text+' More text.'}));
});
