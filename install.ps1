$ErrorActionPreference = "Stop"

$repo = "shishiro26/Kubera"

$release = Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest"
$tag = $release.tag_name
$version = $tag.TrimStart('v')

$arch = "amd64"
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { $arch = "arm64" }

$asset = "kubera_${version}_windows_${arch}.zip"
$url = "https://github.com/$repo/releases/download/$tag/$asset"

$installDir = "$env:LOCALAPPDATA\kubera"
New-Item -ItemType Directory -Force -Path $installDir | Out-Null

$tmpZip = Join-Path $env:TEMP $asset
Write-Host "Downloading kubera $version..."
Invoke-WebRequest -Uri $url -OutFile $tmpZip -UseBasicParsing

Expand-Archive -Path $tmpZip -DestinationPath $installDir -Force
Remove-Item $tmpZip

$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($currentPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$installDir", "User")
}

Write-Host "kubera v$version installed. Restart your terminal and run: kubera"
