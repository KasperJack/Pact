param(
    [string]$ScriptPath,
    [string]$Seed
)


$moduleDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# . blank slate ───────────────────────────────────────────────
$iss = [System.Management.Automation.Runspaces.InitialSessionState]::Create()
$iss.LanguageMode = [System.Management.Automation.PSLanguageMode]::ConstrainedLanguage

# ── 2. load platform in trusted session ──────────────────────────
### Import-Module "$moduleDir/platform.psm1" -Force
## call a constructor from the module using the seed #NO

# . Load & patch the module source ────────────────────────────

$moduleSource = Get-Content "$moduleDir/platform.psm1" -Raw
$moduleSource = $moduleSource -replace 'Builtin-Write-Host', "${seed}-Write-Host"


# ── 2. Load the patched module into a temporary module object ─────
$patchedModule = New-Module -Name "platform_patched" -ScriptBlock ([ScriptBlock]::Create($moduleSource))
$patchedModule = Import-Module $patchedModule -PassThru -Force





# ── 3. inject platform wrapper functions ─────────────────────────
foreach ($name in $patchedModule.ExportedFunctions.Keys) {
    $fn = $patchedModule.ExportedFunctions[$name]
    $iss.Commands.Add((New-Object System.Management.Automation.Runspaces.SessionStateFunctionEntry(
        $name, $fn.ScriptBlock
    )))
}



# ── 4. inject real cmdlets under __real_ names ───────────────────
#    these are callable by your wrappers but not guessable by the script
## loop over some capbilty set by the module for fuct it needs 
$passthroughCmdlets = @(
    @{ Name = "${seed}-Write-Host"; Type = [Microsoft.PowerShell.Commands.WriteHostCommand]  }

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