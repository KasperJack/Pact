#Write-Log "starting in env: $Env"
#Write-Log "tenant: $TenantId"
#__real_Write-Host "asshole"  
#Microsoft.PowerShell.Utility\Write-Host "hacked !!" # static parsing 

#$users = Get-Users -Filter "*"

#foreach ($user in $users) {
#    Write-Log "found user: $($user.Name)"
#}





#[wmiclass]"Win32_Process"
#Invoke-Expression $code







## can't list avallable commands with ##  Get-ChildItem Function: ; Get-Command






#1. Direct Process Spawning
# pass all 
<#
cmd.exe /c calc.exe
cmd /c whoami
[System.Diagnostics.Process]::Start("calc.exe")
Start-Process calc.exe
Start-Process cmd -ArgumentList "/c whoami"
& "C:\Windows\System32\calc.exe"
#>

<# output 
[SANDBOX ERROR] The term 'cmd.exe' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
[SANDBOX ERROR] The term 'cmd.exe' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
[SANDBOX ERROR] Cannot invoke method. Method invocation is supported only on core types in this language mode.
[SANDBOX ERROR] The term 'Start-Process' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
[SANDBOX ERROR] The term 'Start-Process' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
[SANDBOX ERROR] The term 'C:\Windows\System32\calc.exe' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
#>






#2. Expression Evaluation
# pass all 
<#
IEX "cmd /c whoami"
Invoke-Expression "Write-Host bypass"
$cmd = [Convert]::ToBase64String([Text.Encoding]::Unicode.GetBytes("Write-Host pwned"))
powershell -EncodedCommand $cmd
$a = "Write-H"; $b = "ost hi"; IEX ($a + $b)
$x = "IEX"; & $x "Write-Host test"
#>


<#  output 
[SANDBOX ERROR] The term 'Invoke-Expression' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
[SANDBOX ERROR] Cannot invoke method. Method invocation is supported only on core types in this language mode.
[SANDBOX ERROR] The term 'powershell.exe' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
[SANDBOX ERROR] The term 'IEX' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
[SANDBOX ERROR] The term 'IEX' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
#>









#3. .NET Direct Calls
# pass 3/4

<#
[System.Diagnostics.Process]::Start("cmd", "/c whoami")
[Reflection.Assembly]::LoadWithPartialName("System")
[System.Environment]::GetEnvironmentVariables()
[System.Environment]::OSVersion  ## ----> Microsoft Windows NT 10.0.26200.0 
#>


<# output 
[SANDBOX ERROR] Cannot invoke method. Method invocation is supported only on core types in this language mode.
[SANDBOX ERROR] Cannot invoke method. Method invocation is supported only on core types in this language mode.
[SANDBOX ERROR] Cannot invoke method. Method invocation is supported only on core types in this language mode.
#>






#4. Reflection to Access Private/Real Cmdlets
# pass all

<#
$cmd = [Microsoft.PowerShell.Commands.WriteHostCommand]::new()
[AppDomain]::CurrentDomain.GetAssemblies() | % { $_.GetTypes() } 2>$null
$type = [psobject].Assembly.GetType("System.Management.Automation.CommandDiscovery")
$type.GetMethods([Reflection.BindingFlags]"NonPublic,Static") | select Name
#>



<# output 
[SANDBOX ERROR] Cannot create type. Only core types are supported in this language mode.
[SANDBOX ERROR] Cannot invoke method. Method invocation is supported only on core types in this language mode.
[SANDBOX ERROR] Cannot invoke method. Method invocation is supported only on core types in this language mode.
[SANDBOX ERROR] You cannot call a method on a null-valued expression.
#>


















#5. Runspace / Pipeline Escape
# pass all

<#
$rs = [RunspaceFactory]::CreateRunspace()
$rs.Open()
$ps = [PowerShell]::Create()
$ps.Runspace = $rs
$ps.AddScript("whoami").Invoke()
$rs.Close()

# Inline runspace
[PowerShell]::Create().AddScript("Get-Process").Invoke()
#>


<# output 
[SANDBOX ERROR] Cannot invoke method. Method invocation is supported only on core types in this language mode.
[SANDBOX ERROR] You cannot call a method on a null-valued expression.
[SANDBOX ERROR] Cannot invoke method. Method invocation is supported only on core types in this language mode.
[SANDBOX ERROR] The property 'Runspace' cannot be found on this object. Verify that the property exists and can be set.
[SANDBOX ERROR] You cannot call a method on a null-valued expression.
[SANDBOX ERROR] You cannot call a method on a null-valued expression.
[SANDBOX ERROR] Cannot invoke method. Method invocation is supported only on core types in this language mode.
#>








#6. Module / Import Bypass
# pass all 
<#
Import-Module Microsoft.PowerShell.Utility -Force
Remove-Module platform_patched -Force
& (Get-Module Microsoft.PowerShell.Utility) { Write-Host "direct" }
#>


<# output 
[SANDBOX ERROR] The term 'Import-Module' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correctand try again.
[SANDBOX ERROR] The term 'Remove-Module' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correctand try again.
[SANDBOX ERROR] The term 'Get-Module' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
#>







#7. Alias & Command Resolution Tricks
#pass all

<#
Get-Command Write-Host | select *
Set-Alias wh ([Microsoft.PowerShell.Commands.WriteHostCommand])
wh "escaped"
#>

<# output 
[SANDBOX ERROR] The term 'Get-Command' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
[SANDBOX ERROR] The term 'Set-Alias' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct andtry again.
[SANDBOX ERROR] The term 'wh' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
#>











# 8. ScriptBlock & Delegate Tricks
#pass all 

<#
{Write-Host "from block"}.Invoke()

$sb = [scriptblock]::Create("cmd /c whoami")
$sb.Invoke()

Start-Job { Write-Host "in job"; whoami } | Wait-Job | Receive-Job
#>



<# output 
[SANDBOX ERROR] Cannot invoke method. Method invocation is supported only on core types in this language mode.
[SANDBOX ERROR] Cannot invoke method. Method invocation is supported only on core types in this language mode.
[SANDBOX ERROR] You cannot call a method on a null-valued expression.
[SANDBOX ERROR] The term 'Start-Job' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct andtry again.
#>







#9. Environment & File System Probing
#pass all
<#
"test" | Out-File C:\temp\test.txt
Get-Content C:\temp\test.txt

$env:PATH
$env:USERNAME
$env:COMPUTERNAME
dir env:

$env:PATH = "/home" ##done't need env vars inject globals 
#>


<# output 
[SANDBOX ERROR] The term 'Out-File' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
[SANDBOX ERROR] The term 'Get-Content' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
[SANDBOX ERROR] The term 'dir' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
[SANDBOX ERROR] Cannot find drive. A drive with the name 'env' does not exist.
#>






#10. Type Accelerator / COM Objects
#pass all

<#
$shell = New-Object -ComObject WScript.Shell
$shell.Run("calc.exe")
$ie = New-Object -ComObject InternetExplorer.Application
$ie.Navigate("http://example.com")
[activator]::CreateInstance([type]::GetTypeFromProgID("WScript.Shell"))
#>

<# output 
[SANDBOX ERROR] The term 'New-Object' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
[SANDBOX ERROR] You cannot call a method on a null-valued expression.
[SANDBOX ERROR] The term 'New-Object' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a path was included, verify that the path is correct and try again.
[SANDBOX ERROR] You cannot call a method on a null-valued expression.
[SANDBOX ERROR] Cannot invoke method. Method invocation is supported only on core types in this language mode.
#>


[System.Environment]::MachineName
[System.Environment]::UserName
[System.Environment]::CurrentDirectory
[System.Environment]::SystemDirectory
[System.Environment]::Version          # PS/.NET version
[System.Environment]::ProcessorCount
[System.Environment]::Is64BitProcess
[System.Environment]::CommandLine      # how was the sandbox launched?
[System.Environment]::StackTrace 
[System.IO.File]::Exists   # method - will fail
[System.IO.Path]::DirectorySeparatorChar   # property, might work



