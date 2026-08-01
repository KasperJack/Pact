<#
.SYNOPSIS
    Adds this project's bin dirs to PATH for the current PowerShell session.

.DESCRIPTION
    Dot-source this from the repo root (or anywhere) - it resolves paths
    relative to the script's own location, so it works regardless of your
    current working directory. Changes only affect $env:PATH in this
    process; nothing is written to the registry or profile.

.USAGE
    . .\init.ps1
#>

# Resolve the directory this script lives in (repo root), not the CWD
$RepoRoot = $PSScriptRoot

$DirsToAdd = @(
    (Join-Path $RepoRoot 'bin')
    (Join-Path $RepoRoot 'bin\repo')
)

# Split current PATH into a set for quick membership checks (case-insensitive)
$existing = $env:PATH -split ';' | Where-Object { $_ } 
$existingSet = [System.Collections.Generic.HashSet[string]]::new(
    [string[]]$existing,
    [System.StringComparer]::OrdinalIgnoreCase
)

$added = @()
$skippedMissing = @()

foreach ($dir in $DirsToAdd) {
    if (-not (Test-Path $dir)) {
        $skippedMissing += $dir
        continue
    }

    if ($existingSet.Contains($dir)) {
        continue
    }

    # Prepend so these take priority over any same-named tools elsewhere on PATH
    $env:PATH = "$dir;$env:PATH"
    $existingSet.Add($dir) | Out-Null
    $added += $dir
}

if ($added.Count -gt 0) {
    Write-Host "Added to PATH (this session):" -ForegroundColor Green
    $added | ForEach-Object { Write-Host "  $_" -ForegroundColor Green }
} else {
    Write-Host "No new directories added (already on PATH or none existed)." -ForegroundColor Yellow
}

if ($skippedMissing.Count -gt 0) {
    Write-Host "Skipped (does not exist yet):" -ForegroundColor DarkYellow
    $skippedMissing | ForEach-Object { Write-Host "  $_" -ForegroundColor DarkYellow }
}