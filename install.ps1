# Kubera installer for Windows (PowerShell).
# Usage: iwr -useb https://raw.githubusercontent.com/shishiro26/Kubera/main/install.ps1 | iex
#
# Downloads the latest release from GitHub, extracts it to %USERPROFILE%\bin,
# and permanently adds that directory to the User PATH registry key.

$ErrorActionPreference = "Stop"

$Repo       = "shishiro26/Kubera"
$InstallDir = "$env:USERPROFILE\bin"

# ── Fetch latest release ───────────────────────────────────────────────────────
Write-Host "Fetching latest release..."
$release = Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest"
$tag     = $release.tag_name                # e.g. "v1.0.0"
$version = $tag.TrimStart('v')              # e.g. "1.0.0"  (GoReleaser strips v)

# ── Detect architecture ────────────────────────────────────────────────────────
$arch = if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }

$archive = "kubera_${version}_windows_${arch}.zip"
$url     = "https://github.com/$Repo/releases/download/$tag/$archive"
$tmpZip  = "$env:TEMP\kubera_install.zip"
$tmpDir  = "$env:TEMP\kubera_install"

# ── Download ───────────────────────────────────────────────────────────────────
Write-Host "Downloading $archive..."
Invoke-WebRequest -Uri $url -OutFile $tmpZip -UseBasicParsing

# ── Extract ────────────────────────────────────────────────────────────────────
if (Test-Path $tmpDir) { Remove-Item $tmpDir -Recurse -Force }
Expand-Archive -Path $tmpZip -DestinationPath $tmpDir
Remove-Item $tmpZip

# ── Install binary ─────────────────────────────────────────────────────────────
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}
Copy-Item "$tmpDir\kubera.exe" -Destination "$InstallDir\kubera.exe" -Force
Remove-Item $tmpDir -Recurse -Force

Write-Host "Installed kubera $tag to $InstallDir\kubera.exe"

# ── Persist PATH (User scope, registry — no 1024-char setx limit) ─────────────
$current = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($current -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable('Path', "$current;$InstallDir", 'User')
    Write-Host "Added $InstallDir to PATH"
} else {
    Write-Host "$InstallDir is already in PATH"
}

Write-Host ""
Write-Host "Done! Open a new terminal and run: kubera init"
