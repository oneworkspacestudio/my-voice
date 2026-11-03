"""My Voice's torch-free MOSS adapter. OrtCpuRuntime is unmodified upstream.
Reference state is version-specific audio-tokenizer codes, never a portable voice.
Only CPUExecutionProvider is permitted. Inputs are private local files from Go.
"""
import json,sys,wave
from pathlib import Path
import numpy as np
import sentencepiece as spm
from ort_cpu_runtime import OrtCpuRuntime

mode,model,text_or_reference,output,threads=sys.argv[1:6]
runtime=OrtCpuRuntime(model,thread_count=int(threads),max_new_frames=1500,execution_provider='cpu',sample_mode='fixed')
if mode=='clone':
 with wave.open(text_or_reference,'rb') as f:
  assert f.getnchannels()==1 and f.getsampwidth()==2
  rate=f.getframerate();audio=np.frombuffer(f.readframes(f.getnframes()),dtype='<i2').astype(np.float32)/32768
 target=int(runtime.codec_meta['codec_config']['sample_rate'])
 if rate!=target:audio=np.interp(np.arange(round(len(audio)*target/rate))*rate/target,np.arange(len(audio)),audio).astype(np.float32)
 channels=int(runtime.codec_meta['codec_config']['channels'])
 waveform=np.repeat(audio[None,None,:],channels,axis=1)
 results=runtime.sessions['codec_encode'].run(None,{'waveform':waveform,'input_lengths':np.asarray([len(audio)],dtype=np.int32)})
 named=dict(zip([v.name for v in runtime.sessions['codec_encode'].get_outputs()],results))
 length=int(named['audio_code_lengths'].reshape(-1)[0]);codes=np.asarray(named['audio_codes'],dtype=np.int32)[0,:length,:].tolist()
 Path(output).write_text(json.dumps({'engine':'moss','version':'f52645cb','codes':codes}),encoding='utf-8')
else:
 tokenizer=spm.SentencePieceProcessor(model_file=str(runtime.resolve_manifest_relative_path(runtime.manifest['model_files']['tokenizer_model'])))
 if len(sys.argv)>6:
  state=json.loads(Path(sys.argv[6]).read_text(encoding='utf-8'));assert state['engine']=='moss' and state['version']=='f52645cb';codes=state['codes']
 else:raise RuntimeError('A permissioned model-specific voice state is required')
 rows=runtime.build_voice_clone_request_rows(codes,tokenizer.encode(text_or_reference,out_type=int))
 frames=runtime.generate_audio_frames(rows)
 if len(frames)>=1500:raise RuntimeError('Generation reached its duration limit; split the segment before retrying')
 # Streaming decoder bounds peak memory without promising network streaming.
 runtime.codec_streaming_session.reset();chunks=[]
 for n in range(0,len(frames),8):
  result=runtime.codec_streaming_session.run_frames(frames[n:n+8])
  if result is not None:
   audio,length=result
   if length>0:chunks.append(audio[0,:,:length].mean(axis=0))
 if not chunks:raise RuntimeError('Model generated no audio')
 audio=np.concatenate(chunks);pcm=(np.clip(audio,-1,1)*32767).astype('<i2')
 with wave.open(output,'wb') as f:
  f.setnchannels(1);f.setsampwidth(2);f.setframerate(int(runtime.codec_meta['codec_config']['sample_rate']));f.writeframes(pcm.tobytes())
