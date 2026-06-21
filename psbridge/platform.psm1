function Get-Users {
    param([string]$Filter = "*")
    
    return @(
        @{ Name = "Alice"; Email = "alice@company.com" }
        @{ Name = "Bob";   Email = "bob@company.com" }
    )
}

function Write-Log {
    param([string]$Message)
    Microsoft.PowerShell.Utility\Write-Host "[LOG] $Message"
}

function Send-Alert {
    param([string]$Message, [string]$Severity = "info")
    Write-Host "[ALERT][$Severity] $Message"
}


#some error code for security vailations 
function Start-Process { throw "Start-Process is not allowed in this environment"}


function Invoke-Expression {
    $msg = "BLOCKED: Invoke-Expression called at line $($MyInvocation.ScriptLineNumber) in $($MyInvocation.ScriptName)"
    Microsoft.PowerShell.Utility\Write-Host $msg
    throw $msg
}


#`Microsoft.PowerShell.Utility\`,
#`Microsoft.PowerShell.Management\`,
#`Microsoft.PowerShell.Security\`,