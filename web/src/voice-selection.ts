type Voice = {id: string; name: string; engine: string; language: string; preset?: boolean; shared?: boolean};
type Model = {id: string; languages: string[]; shortName?: string};
type Fields = Record<string, any>;

function speechTitle(name: string, voice: Voice | undefined, voices: Voice[]) {
  return voice && voices.some(v => name === 'Voice preview · ' + v.name) ? 'Voice preview · ' + voice.name : name;
}

export function voiceFitsModel(voice: Voice, model: string) {
  return voice.engine === model || (voice.engine === 'pocket' && ['pocket18', 'pocket7'].includes(model));
}

export function speechVoiceOptions(voices: Voice[], models: Model[] = []) {
  return [...voices].sort((a, b) => Number(!!a.preset) - Number(!!b.preset)).map(v => ({
    value: v.id, label: v.name + (v.preset ? '' : (v.shared ? ' · Shared' : ' · Clone')+(models.find(m=>m.id===v.engine)?.shortName?' · '+models.find(m=>m.id===v.engine)!.shortName:''))
  }));
}

export function selectSpeechVoice(fields: Fields, id: string, voices: Voice[], models: Model[]) {
  const voice = voices.find(v => v.id === id);
  if (!voice) return fields;
  const model = voice.engine === 'pocket' ? (['pocket18', 'pocket7'].includes(fields.model) ? fields.model : 'pocket18') : voice.engine;
  const languages = models.find(m => m.id === model)?.languages || [voice.language];
  const language = languages.includes(fields.language) ? fields.language : languages.includes(voice.language) ? voice.language : languages[0];
  // A preview opened as a draft has an automatic title. Keep custom titles and
  // historical recordings intact, but do not name a new clone generation Alba.
  const projectName = speechTitle(fields.projectName, voice, voices);
  return {...fields, voice: voice.id, model, language, projectName};
}

export function selectSpeechModel(fields: Fields, model: string, voices: Voice[], models: Model[]) {
  const matching = voices.filter(v => voiceFitsModel(v, model));
  const voice = matching.find(v => v.id === fields.voice) || matching.find(v => v.preset) || matching[0];
  return voice ? selectSpeechVoice({...fields, model}, voice.id, voices, models) : {...fields, model};
}

export function recordingVoiceLabel(script: any, voices: Voice[], models: Model[] = []) {
  if (!script) return '';
  const ids: string[] = Array.from(new Set(script.inputMode === 'segments' ? script.segments.map((s: any) => s.voice) : [script.voice]));
  const names = ids.map(id => voices.find(v => v.id === id)?.name || 'Unavailable voice');
  return (names.length > 1 ? 'Voices · ' : 'Voice · ') + names.join(', ') + (models.find(m=>m.id===script.model)?.shortName?' · '+models.find(m=>m.id===script.model)!.shortName:'');
}

export function recordingTitle(script: any, voices: Voice[]) {
  return script ? speechTitle(script.name, voices.find(v => v.id === script.voice), voices) : '';
}
