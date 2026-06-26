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

 
    $raw = '{"fid":1,"params":{"Year":2022,"Month":12}}'


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


