param([Parameter(Mandatory=$true)][string]$Headers,[Parameter(Mandatory=$true)][string]$Output)
$ErrorActionPreference='Stop'
$vswhere=Join-Path ${env:ProgramFiles(x86)} 'Microsoft Visual Studio/Installer/vswhere.exe'
$installation=& $vswhere -latest -products '*' -requires Microsoft.VisualStudio.Component.VC.Tools.x86.x64 -property installationPath
if (!$installation) { throw 'Install Visual Studio C++ build tools to rebuild the resident Qwen adapter.' }
$vcvars=Join-Path $installation 'VC/Auxiliary/Build/vcvars64.bat'
$source=Join-Path $PSScriptRoot '../runtime-src/qwen_resident.cpp'
$build=Join-Path (Split-Path $Headers) 'adapter-build'
New-Item -ItemType Directory -Force $build | Out-Null
$object=Join-Path $build 'qwen_resident.obj'
$library=Join-Path $build 'myvoice-qwen.lib'
$command=Join-Path $build 'compile.cmd'
@("@echo off", "call `"$vcvars`" >nul", 'if errorlevel 1 exit /b 1', "cl /nologo /std:c++17 /O2 /MD /EHsc /Brepro /LD /I `"$Headers`" `"$source`" /Fo:`"$object`" /Fe:`"$Output`" /link /BREPRO /IMPLIB:`"$library`"") | Set-Content -LiteralPath $command -Encoding ascii
& cmd /d /c $command
if ($LASTEXITCODE -ne 0) { throw 'Resident Qwen adapter compilation failed' }
