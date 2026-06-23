function Get-Users {
    param([string]$Filter = "*")
    return @(
        [PSCustomObject]@{ Name = "Alice"; Email = "alice@company.com" }
        [PSCustomObject]@{ Name = "Bob";   Email = "bob@company.com"   }
    )
}

function Write-Log {
    param([string]$Message)
   
    Builtin-Write-Host "[LOG] $Message"
}



function Send-Alert {
    param([string]$Message, [string]$Severity = "info")
    Builtin-Write-Host "[ALERT][$Severity] $Message"
}



function Private-Somthing {
    
    Builtin-Write-Host "DOING SOMTHING PRIVATE"
}



function Write-Host {
    param(
        [Parameter(ValueFromPipeline)][object]$Object,
        [string]$ForegroundColor,
        [switch]$NoNewline
    )

     Builtin-Write-Host "gg"   
}


#Export-ModuleMember -Function 'Write-Log'