
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8




Import-Module Microsoft.PowerShell.Management

# Pre-instantiate COM objects you'll reuse
$Global:Shell = New-Object -ComObject WScript.Shell

# Pre-load .NET assemblies
Add-Type -AssemblyName System.Windows.Forms





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


$Router = @{}

$Router[0] = {
    param($req)

    # warm up execution 
    $raw = '{"fid":1,"params":{"Year":2022,"Month":12}}'


    $req = $raw | ConvertFrom-Json

    $handler = $Router[1]

    & $handler $req | ConvertTo-Json -Compress| Out-Null
    $null = $Global:Shell.CreateShortcut("C:\warmup.lnk")


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



$Router[2] = {
    param($req)

    $name      = $req.params.Name        # "My App"
    $target    = $req.params.Target      # "C:\Program Files\MyApp\app.exe"
    $icon      = $req.params.Icon        # optional, can be $null
    $argsf      = $req.params.Arguments   # optional, can be $null

    $desktop = [System.Environment]::GetFolderPath("Desktop")
    $path    = Join-Path $desktop "$name.lnk"

    #$shell    = New-Object -ComObject WScript.Shell
    $shortcut = $Global:Shell.CreateShortcut($path)

    $shortcut.TargetPath = $target

    if ($argsf)  { $shortcut.Arguments       = $argsf }
    if ($icon)  { $shortcut.IconLocation    = $icon }

    $shortcut.Save()

    return @{ ok = $true; result = $path }
}













while ($true) {
    $raw = [Console]::In.ReadLine()
    if ($null -eq $raw) { break }

    try {
        $req = $raw | ConvertFrom-Json
    }
    catch {
        [Console]::WriteLine((@{ ok = $false; error = "invalid json" } | ConvertTo-Json -Compress))
        continue
    }


    $fid = [int]$req.fid

    $handler = $Router[$fid]

    if ($null -ne $handler) {



        try {
            $res = & $handler $req
        }
        catch {
            [Console]::WriteLine((@{ ok = $false; error = $_.Exception.Message } | ConvertTo-Json -Compress))
            continue
        }




        [Console]::WriteLine(($res | ConvertTo-Json -Compress))






    } else {

       [Console]::WriteLine( (@{ ok = $false; error = "unknown fid: $($fid)" } | ConvertTo-Json -Compress))
    }
}








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