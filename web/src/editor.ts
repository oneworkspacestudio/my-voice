type Fields = Record<string, any>;
type Segment = {id: string; text: string; voice: string; [key: string]: any};

export function splitEditorText(text: string, voice: string): Segment[] {
  const encoder = new TextEncoder();
  const result: Segment[] = [];
  for (const sentence of text.match(/[^.!?\n]+[.!?]*|[.!?]+/gu) || []) {
    let chunk = '', bytes = 0;
    for (const char of sentence.trim()) {
      const size = encoder.encode(char).length;
      if (bytes + size > 500) {
        const space = chunk.lastIndexOf(' ');
        const split = space > 0 && encoder.encode(chunk.slice(0, space)).length > 150 ? space : chunk.length;
        if (chunk.slice(0, split).trim()) result.push({id: crypto.randomUUID(), text: chunk.slice(0, split).trim(), voice});
        chunk = chunk.slice(split).trim(); bytes = encoder.encode(chunk).length;
      }
      chunk += char; bytes += size;
    }
    if (chunk.trim()) result.push({id: crypto.randomUUID(), text: chunk.trim(), voice});
  }
  return result;
}

export function editedSegments(fields: Fields, segments: Segment[]): Segment[] {
  return segments.map(s => ({...s, text: fields['text_' + s.id] ?? s.text, voice: fields['voice_' + s.id] ?? s.voice}));
}

export function buildScript(fields: Fields, segments: Segment[], segmentMode: boolean, dictionary: any[]) {
  const edited = segmentMode ? editedSegments(fields, segments) : [];
  return {name: fields.projectName, text: segmentMode ? edited.map(s => s.text).join('\n') : fields.script,
    inputMode: segmentMode ? 'segments' : 'text', voice: fields.voice, model: fields.model,
    language: fields.language, speed: Number(fields.speed), pitch: Number(fields.pitch),
    volume: Number(fields.volume), pause: Number(fields.pause), dictionary, segments: edited};
}

export function loadDraft(script: any) {
  let segments: Segment[] = script.segments || [];
  // Legacy auto-chunked recordings are full-text drafts. Preserve explicitly
  // selected segment mode and legacy multi-voice/segment-only projects.
  const segmentMode = script.inputMode === 'segments' || (!script.inputMode &&
    (segments.some(s => s.voice !== script.voice) || (!script.text?.trim() && segments.length > 0)));
  if (!segmentMode && segments.map(s => s.text).join('').replace(/\s/gu, '') !== (script.text || '').replace(/\s/gu, '')) segments = [];
  const fields: Fields = {script: segmentMode ? segments.map(s => s.text).join('\n') : script.text || '',
    projectName: script.name, voice: script.voice, model: script.model, language: script.language,
    speed: script.speed ?? 1, pitch: script.pitch ?? 0, volume: script.volume ?? 1, pause: script.pause ?? .3,
    rules: (script.dictionary || []).map((r: any) => r.from + '=' + r.to).join('\n'),
    ...segmentFields(segments)};
  return {fields, segments, segmentMode};
}

export function segmentFields(segments: Segment[]) {
  return Object.fromEntries(segments.flatMap(s => [['text_' + s.id, s.text], ['voice_' + s.id, s.voice]]));
}

export function scriptSignature(script: ReturnType<typeof buildScript>) {
  return JSON.stringify([script.name, script.text, script.inputMode, script.voice, script.model,
    script.language, script.speed, script.pitch, script.volume, script.pause, script.dictionary,
    script.segments.map(s => [s.id, s.text, s.voice])]);
}
