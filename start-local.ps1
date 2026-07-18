#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Starts Trickreport backend and frontend locally without Docker.
.DESCRIPTION
    Creates the database (if it does not exist), applies migrations, starts the
    Go backend and the Astro frontend dev server.
.PARAMETER PostgresPassword
    Password for the local postgres user.
.PARAMETER SkipDb
    Skip database creation and migrations.
.PARAMETER SkipBackend
    Skip starting the backend.
.PARAMETER SkipFrontend
    Skip starting the frontend.
#>
param(
    [Parameter(Mandatory = $true)]
    [string]$PostgresPassword,

    [switch]$SkipDb,
    [switch]$SkipBackend,
    [switch]$SkipFrontend
)

$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $MyInvocation.MyCommand.Definition

$env:PGPASSWORD = $PostgresPassword

function Invoke-PSQL($arguments) {
    $output = & psql @arguments 2>&1
    if ($LASTEXITCODE -ne 0) {
        throw "psql failed: $output"
    }
    return $output
}

if (-not $SkipDb) {
    Write-Host 'Ensuring database exists...' -ForegroundColor Cyan
    $exists = Invoke-PSQL @('-U', 'postgres', '-tAc', 'SELECT 1 FROM pg_database WHERE datname=''trickreport''')
    if (-not $exists -or $exists.Trim() -ne '1') {
        Invoke-PSQL @('-U', 'postgres', '-c', 'CREATE DATABASE trickreport;')
    } else {
        Write-Host '  Database already exists.'
    }

    Write-Host 'Applying migrations...' -ForegroundColor Cyan
    $migrations = @(
        '001_init.sql',
        '002_tickets.sql',
        '003_knowledge_base.sql',
        '004_sla_policies.sql',
        '005_automations.sql',
        '006_sla_breached.sql'
    )
    foreach ($m in $migrations) {
        $path = Join-Path $root 'backend' 'migrations' $m
        Invoke-PSQL @('-U', 'postgres', '-d', 'trickreport', '-f', $path)
        Write-Host "  Applied $m"
    }
}

$backendJob = $null
$frontendJob = $null

try {
    if (-not $SkipBackend) {
        $backendPath = Join-Path $root 'backend'
        $backendEnv = Join-Path $backendPath '.env'
        if (-not (Test-Path $backendEnv)) {
            Copy-Item (Join-Path $backendPath '.env.example') $backendEnv
            Write-Host "Created $backendEnv. Please review it before running again if you need custom values." -ForegroundColor Yellow
        }

        Write-Host 'Generating Templ files...' -ForegroundColor Cyan
        Push-Location $backendPath
        templ generate ./internal/interfaces/http/views/
        Pop-Location

        Write-Host 'Starting Go backend on http://localhost:8080 ...' -ForegroundColor Cyan
        $backendJob = Start-Job -ScriptBlock {
            param($path)
            Set-Location $path
            go run cmd/api/main.go
        } -ArgumentList $backendPath
    }

    if (-not $SkipFrontend) {
        $frontendPath = Join-Path $root 'frontend'
        $frontendEnv = Join-Path $frontendPath '.env'
        if (-not (Test-Path $frontendEnv)) {
            Copy-Item (Join-Path $frontendPath '.env.example') $frontendEnv
            Write-Host "Created $frontendEnv. Please review it before running again if you need custom values." -ForegroundColor Yellow
        }

        Write-Host 'Installing frontend dependencies...' -ForegroundColor Cyan
        Push-Location $frontendPath
        npm install
        Pop-Location

        Write-Host 'Starting Astro frontend on http://localhost:4321 ...' -ForegroundColor Cyan
        $frontendJob = Start-Job -ScriptBlock {
            param($path)
            Set-Location $path
            npm run dev
        } -ArgumentList $frontendPath
    }

    Write-Host ''
    Write-Host 'Both services are starting in background jobs.' -ForegroundColor Green
    Write-Host 'Press Ctrl+C to stop.' -ForegroundColor Green

    while ($true) {
        if ($backendJob -and $backendJob.State -eq 'Failed') {
            throw "Backend failed: $(Receive-Job -Job $backendJob)"
        }
        if ($frontendJob -and $frontendJob.State -eq 'Failed') {
            throw "Frontend failed: $(Receive-Job -Job $frontendJob)"
        }
        Start-Sleep -Seconds 1
    }
}
finally {
    if ($backendJob) {
        Stop-Job -Job $backendJob -ErrorAction SilentlyContinue
        Remove-Job -Job $backendJob -ErrorAction SilentlyContinue
    }
    if ($frontendJob) {
        Stop-Job -Job $frontendJob -ErrorAction SilentlyContinue
        Remove-Job -Job $frontendJob -ErrorAction SilentlyContinue
    }
}
