export type CloneModel = {id:string;name:string;cloning:boolean;cloningReady?:boolean;installed:boolean;installing?:boolean;status?:string;size:number;languages:string[];summary?:string;ramGB?:number};

export const modelInstalling = (m:CloneModel) => !!m.installing || /^(Downloading|Preparing)/.test(m.status||'');
export const cloneReady = (m?:CloneModel) => !!m && m.installed && (m.cloningReady ?? m.cloning) && !modelInstalling(m);

export function cloningSetup(models:CloneModel[], requested:string, loaded:boolean, error:string, owner:boolean){
 const choices=models.filter(m=>m.cloning||m.cloningReady);
 const model=choices.find(m=>m.id===requested)||choices.find(cloneReady)||choices[0];
 const busy=!!model&&modelInstalling(model),ready=cloneReady(model)&&!error&&loaded;
 const failed=!!model?.status&&!busy&&!ready&&model.status!=='Ready';
 return {model,ready,busy,options:choices.map(m=>({value:m.id,label:m.name})),
  title:!loaded?'Checking cloning models…':error?'Could not check cloning models':!model?'No cloning model is available':ready?'Ready to record':busy?'Getting your cloning model ready':failed?'Model setup needs attention':'Install a model before recording',
  detail:!loaded?'Checking this PC before you spend time recording.':error?'Refresh the model list to try again. Your reference and voice details are kept.':!model?'Refresh the model list or update My Voice through Workspace Updates.':ready?model.name+' is installed and ready on this PC.':busy?(model.status||'Preparing model…'):failed?model.status:model.name+' · '+Math.ceil(model.size/1024/1024)+' MB download · runs locally on your CPU.',
  status:ready?'Ready':busy?'Installing':failed?'Setup failed':'Not ready',
  canInstall:loaded&&!error&&!!model&&!ready&&!busy&&owner,
  installLabel:failed?'Retry model setup':'Install cloning model',
  needsOwner:loaded&&!error&&!!model&&!ready&&!owner,
  progress:busy&&model?.status?.startsWith('Downloading')?Number(model.status.match(/(\d+)%/)?.[1]||0):undefined
 };
}

export async function requireCloningModel(load:()=>Promise<CloneModel[]>, requested:string){
 const models=await load();
 const model=models.find(m=>m.id===requested);
 if(!cloneReady(model))throw Error('Install the selected cloning model on this screen and wait until it is ready before adding a reference.');
 return model!;
}
