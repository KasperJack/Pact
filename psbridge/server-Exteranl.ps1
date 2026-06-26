
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

#$stdin = [System.IO.StreamReader]::new([Console]::OpenStandardInput(), [System.Text.Encoding]::UTF8, $true, 65536)

Add-Type -Path ".\Newtonsoft.Json.13.0.4\lib\net45\Newtonsoft.Json.dll"


function Read-Json($raw) {
    [Newtonsoft.Json.JsonConvert]::DeserializeObject(
        $raw,
        [System.Collections.Generic.Dictionary[string,object]]
    )
}

function Write-Json($obj) {
    [Newtonsoft.Json.JsonConvert]::SerializeObject($obj)
}




$Router = @{}

$Router[0] = {
    param($req)

    
    $raw = '{"fid":1,"params":{"Year":2022,"Month":12}}'


    $req = Read-Json $raw

    $handler = $Router[1]

    & $handler $req | Out-Null
    Write-Json @{ ok = $true; result = "warm" } | Out-Null


    return @{ ok = $true; result = "ready" }
    
}





$Router[1] = {
    param($req)

    $year  = [int]$req["params"]["Year"]
    $month = [int]$req["params"]["Month"]

    $result = (Get-Date -Year $year -Month $month).ToString("o")

    return @{
        ok     = $true
        result = $result
    }
}




while ($true) {
    $raw = [Console]::In.ReadLine()
    if ($null -eq $raw) { break }

    try {
        $req = Read-Json $raw
    }
    catch {
        [Console]::WriteLine((Write-Json @{ ok = $false; error = $_.Exception.Message }))
        continue   
    }

    $fid = [int]$req["fid"]        
    $handler = $Router[$fid]

    if ($null -ne $handler) {
        $res = & $handler $req
        [Console]::WriteLine((Write-Json $res))
    } else {
        [Console]::WriteLine((Write-Json @{ ok = $false; error = "unknown action: $fid" }))
    }
}

































function Init-Server {
    #  pipeline system
    #"warmup" | Out-Null

  
    Get-Date | Out-Null
    Get-Process | Out-Null

   
    @{a=1} | ConvertTo-Json | Out-Null
}

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








Init-Server
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








<#
----
doing 7 rqs
server strated in: 361.721ms
first rquest: 646µs
elapsed (server): 5.3467ms



--------

doing 50000 rqs
server strated in: 375.6601ms
first rquest: 1.0688ms
elapsed (server): 31.4554707s



-------
doing 10000 rqs
server strated in: 418.3642ms
first rquest: 1.6245ms
elapsed (server): 6.0196031s














q/ Runspaces concurrency ? 

#>