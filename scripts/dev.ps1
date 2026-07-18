# Trickreport local dev script (PowerShell)
# Starts backend and frontend concurrently. Press Ctrl+C to stop both.

$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot

Write-Host "Starting Trickreport dev environment (local, no Docker)..." -ForegroundColor Cyan
Write-Host ""

# Check prerequisites
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: Go is not installed or not in PATH" -ForegroundColor Red
    exit 1
}
if (-not (Get-Command npm -ErrorAction SilentlyContinue)) {
    Write-Host "ERROR: npm is not installed or not in PATH" -ForegroundColor Red
    exit 1
}

# Check backend .env
if (-not (Test-Path "$root\backend\.env")) {
    Write-Host "WARNING: backend/.env not found. Copying from .env.example" -ForegroundColor Yellow
    Copy-Item "$root\backend\.env.example" "$root\backend\.env"
    Write-Host "Edit backend/.env with your PostgreSQL credentials before running again." -ForegroundColor Yellow
    exit 1
}

# Check frontend node_modules
if (-not (Test-Path "$root\frontend\node_modules")) {
    Write-Host "Installing frontend dependencies..." -ForegroundColor Cyan
    Push-Location "$root\frontend"
    npm install
    Pop-Location
}

# Start backend
Write-Host "[backend] Starting on http://localhost:8080 ..." -ForegroundColor Green
$backend = Start-Process -FilePath "go" -ArgumentList "run", "cmd/api/main.go" -WorkingDirectory "$root\backend" -PassThru -NoNewWindow

# Start frontend
Write-Host "[frontend] Starting on http://localhost:4321 ..." -ForegroundColor Green
$frontend = Start-Process -FilePath "npm" -ArgumentList "run", "dev" -WorkingDirectory "$root\frontend" -PassThru -NoNewWindow

Write-Host ""
Write-Host "Both services are running. Press Ctrl+C to stop." -ForegroundColor Cyan
Write-Host ""

try {
    while ($true) {
        Start-Sleep -Seconds 1
    }
} finally {
    Write-Host ""
    Write-Host "Stopping services..." -ForegroundColor Yellow
    if ($backend -and -not $backend.HasExited) { Stop-Process -Id $backend.Id -Force }
    if ($frontend -and -not $frontend.HasExited) { Stop-Process -Id $frontend.Id -Force }
    Write-Host "Stopped." -ForegroundColor Green
}
