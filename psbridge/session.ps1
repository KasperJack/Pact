param([string]$ScriptPath)

$moduleDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# ── 1. blank slate ───────────────────────────────────────────────
$iss = [System.Management.Automation.Runspaces.InitialSessionState]::Create()
$iss.LanguageMode = [System.Management.Automation.PSLanguageMode]::ConstrainedLanguage

# ── 2. load platform in trusted session ──────────────────────────
Import-Module "$moduleDir/platform.psm1" -Force

# ── 3. inject platform wrapper functions ─────────────────────────
foreach ($name in @('Get-Users', 'Write-Log', 'Send-Alert', 'Write-Host', 'Copy-Item')) { ##from platfrom module 
    $fn = Get-Item "Function:\$name"
    $iss.Commands.Add((New-Object System.Management.Automation.Runspaces.SessionStateFunctionEntry(
        $name, $fn.ScriptBlock
    )))
}

# ── 4. inject real cmdlets under __real_ names ───────────────────
#    these are callable by your wrappers but not guessable by the script
$passthroughCmdlets = @(
    @{ Name = '__real_Write-Host'; Type = [Microsoft.PowerShell.Commands.WriteHostCommand]  }
    @{ Name = '__real_Copy-Item';  Type = [Microsoft.PowerShell.Commands.CopyItemCommand]   }
)

foreach ($c in $passthroughCmdlets) {
    $iss.Commands.Add((New-Object System.Management.Automation.Runspaces.SessionStateCmdletEntry(
        $c.Name, $c.Type, $null
    )))
}

# ── 5. inject globals ────────────────────────────────────────────
foreach ($pair in @(
    @{ Name = 'Env';        Value = 'production' }
    @{ Name = 'TenantId';   Value = 'abc-123'    }
    @{ Name = 'MaxRetries'; Value = 3             }
)) {
    $iss.Variables.Add((New-Object System.Management.Automation.Runspaces.SessionStateVariableEntry(
        $pair.Name, $pair.Value, ""
    )))
}

# ── 6. run ───────────────────────────────────────────────────────
$code = Get-Content $ScriptPath -Raw
$runspace = [System.Management.Automation.Runspaces.RunspaceFactory]::CreateRunspace($iss)
$runspace.Open()

try {
    $ps = [System.Management.Automation.PowerShell]::Create()
    $ps.Runspace = $runspace
    $ps.AddScript($code) | Out-Null
    $ps.Invoke() | ForEach-Object { Microsoft.PowerShell.Utility\Write-Host $_ }

    $ps.Streams.Information | ForEach-Object { 
        Microsoft.PowerShell.Utility\Write-Host $_.MessageData 
    }
    $ps.Streams.Error | ForEach-Object {
        Microsoft.PowerShell.Utility\Write-Host "[SANDBOX ERROR] $_" -ForegroundColor Red
    }
} finally {
    $runspace.Close()
    $runspace.Dispose()
}