function Get-Users {
    param([string]$Filter = "*")
    return @(
        [PSCustomObject]@{ Name = "Alice"; Email = "alice@company.com" }
        [PSCustomObject]@{ Name = "Bob";   Email = "bob@company.com"   }
    )
}

function Write-Log {
    param([string]$Message)
    __real_Write-Host "[LOG] $Message"
}

function Send-Alert {
    param([string]$Message, [string]$Severity = "info")
    __real_Write-Host "[ALERT][$Severity] $Message"
}


function Write-Host {
    param(
        [Parameter(ValueFromPipeline)][object]$Object,
        [string]$ForegroundColor,
        [switch]$NoNewline
    )

     __real_Write-Host "hole agin"   
    # your checks/logging here
    __real_Write-Host $Object
}

function Copy-Item {
    param([string]$Path, [string]$Destination)
    # your checks here
    if ($Path -notmatch '^C:\\allowed\\') {
        __real_Write-Host "[BLOCKED] Copy-Item: $Path"
        return
    }
    __real_Copy-Item $Path $Destination
}