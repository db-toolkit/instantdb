# InstantDB Windows Installer
$ErrorActionPreference = "Stop"

Write-Host "🚀 Installing instant-db..." -ForegroundColor Cyan

# Detect architecture
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

# Get latest release
Write-Host "📦 Fetching latest release..." -ForegroundColor Cyan
$release = Invoke-RestMethod -Uri "https://api.github.com/repos/db-toolkit/instantdb/releases/latest"
$version = $release.tag_name

$binaryName = "instant-db-windows-$arch.exe"
$downloadUrl = "https://github.com/db-toolkit/instantdb/releases/download/$version/$binaryName"

Write-Host "📥 Downloading instant-db $version for windows-$arch..." -ForegroundColor Cyan
$tempFile = "$env:TEMP\instant-db.exe"

try {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile
} catch {
    Write-Host "❌ Download failed: $_" -ForegroundColor Red
    exit 1
}

# Install to user's local bin
$installDir = "$env:LOCALAPPDATA\Programs\instant-db"
New-Item -ItemType Directory -Force -Path $installDir | Out-Null
Move-Item -Path $tempFile -Destination "$installDir\instant-db.exe" -Force

# Add to PATH if not already there
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
    Write-Host "✅ Added to PATH (restart terminal to use)" -ForegroundColor Green
}

Write-Host ""
Write-Host "✅ instant-db $version installed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "🎉 Get started:" -ForegroundColor Cyan
Write-Host "   instant-db start" -ForegroundColor White
Write-Host ""
Write-Host "📖 For more info: https://github.com/db-toolkit/instant-db" -ForegroundColor Cyan
