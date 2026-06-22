Write-Log "starting in env: $Env"
Write-Log "tenant: $TenantId"
__real_Write-Host "asshole"  
#Microsoft.PowerShell.Utility\Write-Host "hacked !!" # static parsing 

$users = Get-Users -Filter "*"

foreach ($user in $users) {
    Write-Log "found user: $($user.Name)"
}

Send-Alert -Message "script finished" -Severity "info"




#[wmiclass]"Win32_Process"

#Invoke-Expression $code