# sidecar.ps1 // RPC
$stdin = [Console]::OpenStandardInput()



Add-Type -Path ".\Google.Protobuf.3.35.1\lib\net5.0\Google.Protobuf.dll"

$ref = [Google.Protobuf.MessageParser].Assembly.Location



$runtimeDir = [System.Runtime.InteropServices.RuntimeEnvironment]::GetRuntimeDirectory()
$frameworkDir = Split-Path ([System.Object].Assembly.Location)


$refs = @(
    $ref,
    (Join-Path $frameworkDir "System.Runtime.dll"),
    (Join-Path $frameworkDir "System.Collections.dll"),
    (Join-Path $frameworkDir "netstandard.dll")
)

Add-Type `
    -Path ".\m.cs" `
    -ReferencedAssemblies $refs `




$Router = @{}

$Router[0] = {
    param($req)

    # warm real execution path
    $raw = '{"fid":1,"params":{"Year":2022,"Month":12}}'

    # 🔥
    $req = $raw | ConvertFrom-Json

    $handler = $Router[1]

    & $handler $req | Out-Null





    return @{ ok = $true; result = "ready" }
    
}


$Router[1] = {
    param($req)

    $params = ConvertTo-Hashtable $req.params

    $result = Get-Date @params

    return @{
        ok     = $true
        result = $result.ToString("o")
    }
}




while ($true) {

    $lenBuf = New-Object byte[] 4
    $fidBuff = New-Object byte[] 2


    $n = $stdin.Read($lenBuf,0,4)
    $f = $stdin.Read($lenBuf,0,2)

    if ($n -eq 0) {
        break
    }

    $length = [BitConverter]::ToUInt32($lenBuf,0)

    payloadBuff = New-Object byte[] $length


    while ($read -lt $length) {
        $read += $stdin.Read(
            $payloadBuff,
            $read,
            $length-$read
        )
    }




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