param([string]$Go='go',[string]$Version='0.1.0-preview.8',[string]$Out='build/MyVoice.exe')
$ErrorActionPreference='Stop'
$root=Split-Path $PSScriptRoot -Parent
if(!(Test-Path -LiteralPath (Join-Path $root 'internal/app/payload.zip'))){throw 'Run scripts/prepare-runtime.py with a project-owned cache folder first.'}
Push-Location $root
try{
 & npm --prefix web ci --no-audit --no-fund
 if($LASTEXITCODE){throw 'Frontend dependencies failed'}
 & npm --prefix web run build
 if($LASTEXITCODE){throw 'Frontend build failed'}
 $env:CGO_ENABLED='0'
 & $Go build -buildvcs=false -trimpath -ldflags "-s -w -X main.version=$Version" -o $Out ./cmd/my-voice
 if($LASTEXITCODE){throw 'Native Windows build failed'}
}finally{Pop-Location}
