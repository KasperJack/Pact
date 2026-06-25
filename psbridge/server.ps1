# sidecar.ps1 // RPC
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8




function ConvertTo-Hashtable($obj) {
    $h = @{}

    if ($null -eq $obj) {
        return $h
    }

    foreach ($p in $obj.PSObject.Properties) {
        $h[$p.Name] = $p.Value
    }

    return $h
}






function Init-Server {
    #  pipeline system
    #"warmup" | Out-Null

  
    Get-Date | Out-Null
    Get-Process | Out-Null

   
    @{a=1} | ConvertTo-Json | Out-Null
}

Init-Serve

while ($true) {
    #write-host "starting loopp .."
    $raw = [Console]::In.ReadLine()

    #Write-Host "RAW=[$raw]"


    if ($null -eq $raw) { break }  # stdin closed = exit

    try {
      
        $req = $raw | ConvertFrom-Json

        $res = switch ($req.fid) {


            1{#Get-Date
                $params = ConvertTo-Hashtable $req.params

                @{
                ok = $true
                result =  (Get-Date @params).ToString("o")
                }

            }



            "delete-reg-key" {
                Remove-Item -Path $req.path -Recurse -Force | Out-Null
                @{ ok = $true }
            }
            "set-reg-value" {
                Set-ItemProperty -Path $req.path -Name $req.name -Value $req.value | Out-Null
                @{ ok = $true }
            }
            "add-to-path" {
                $current = [Environment]::GetEnvironmentVariable("Path", "Machine")
                [Environment]::SetEnvironmentVariable("Path", "$current;$($req.value)", "Machine")
                @{ ok = $true }
            }
            "create-shortcut" {
                $shell = New-Object -ComObject WScript.Shell
                $s = $shell.CreateShortcut($req.lnkPath)
                $s.TargetPath = $req.targetPath
                $s.Save()
                @{ ok = $true }
            }
            default {
                @{ ok = $false; error = "unknown fid: $($req.fid)" }
            }
        }
    }
    catch {
        $res = @{ ok = $false; error = $_.Exception.Message }

    }

    $res | ConvertTo-Json -Compress
}




#{"fid":1,"params":{}}




<#


In your case it’s probably one (or more) of these:

no clear input contract (schema not enforced)
no error taxonomy (everything is just "error": "...")
no edge case pressure (null params, invalid JSON, empty stdin)
no reproducible test cases
no separation between:
transport (stdin/stdout)
protocol (JSON shape)
execution (PowerShell actions)




Try:

missing params
invalid JSON
unknown fid
null stdin
huge payload
wrong types



stdin → parser → router → executor → serializer → stdout
#>



<#

Go
↓
pipe
↓
PowerShell
↓
pipe
↓
Go
#>







<#

1000 → server 1.01s | batch 1.12s
2000 → server 1.56s | batch 1.89s
3000 → server 2.10s | batch 2.67s
5000 → server 3.35s | batch 4.17s
6000 → server 3.70s | batch 4.95s
7000 → server 4.16s | batch 5.69s


doing 200 rqs
server strated in: 12.0716ms
first rquest: 406.0983ms // cold start 
slepping for 30s
elapsed (server): 195.0404ms

#>
