#Requires -Version 5.1
[CmdletBinding()]
param()

$ErrorActionPreference = 'Stop'

$Repo = 'WarnetBes/cursor-tool'
$Binary = 'cursor-tool'
$InstallDir = "$env:LOCALAPPDATA\cursor-tool"

function Get-LatestRelease {
    $apiUrl = "https://api.github.com/repos/$Repo/releases/latest"
    $response = Invoke-RestMethod -Uri $apiUrl -Headers @{ 'User-Agent' = 'cursor-tool-installer' }
    return $response.tag_name
}

function Get-Arch {
    if ([Environment]::Is64BitOperatingSystem) {
        return 'amd64'
    }
    return '386'
}

Write-Host "cursor-tool installer for Windows" -ForegroundColor Cyan
Write-Host ""

$version = Get-LatestRelease
$arch = Get-Arch
$filename = "${Binary}_windows_${arch}.zip"
$downloadUrl = "https://github.com/$Repo/releases/download/$version/$filename"

Write-Host "Installing $Binary $version for windows/$arch..."

if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

$tmpFile = Join-Path $env:TEMP $filename
Invoke-WebRequest -Uri $downloadUrl -OutFile $tmpFile -UseBasicParsing

$extractDir = Join-Path $env:TEMP "cursor-tool-extract"
if (Test-Path $extractDir) { Remove-Item $extractDir -Recurse -Force }
Expand-Archive -Path $tmpFile -DestinationPath $extractDir -Force

$exePath = Join-Path $extractDir "${Binary}.exe"
if (-not (Test-Path $exePath)) {
    throw "Binary not found in archive."
}

Copy-Item $exePath (Join-Path $InstallDir "${Binary}.exe") -Force
Remove-Item $tmpFile -Force
Remove-Item $extractDir -Recurse -Force

$currentPath = [Environment]::GetEnvironmentVariable('PATH', 'User')
if ($currentPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable('PATH', "$currentPath;$InstallDir", 'User')
    Write-Host "Added $InstallDir to PATH"
}

Write-Host "$Binary installed to $InstallDir" -ForegroundColor Green
Write-Host "Run '$Binary --version' in a new terminal to verify installation."
