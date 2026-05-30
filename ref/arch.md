manifest.lua  →  gopher-lua  →  Go functions  →  OS
   (portable)      (bridge)      (API)     (specific ps/sh)


//
verify = function(ctx)
    ctx.assert_file(ctx.dir .. "/python.exe")
    ctx.assert_shim("python")
    ctx.assert_runs("python --version")
end


or maybe do assersion from the go side 
Wrap execution in a goroutine with a timeout ?? 




 the cleanest way to do it:
goL := lua.NewState(lua.Options{SkipOpenLibs: true}) // open NOTHING
defer L.Close()
Then manually expose only what you want:
go// expose only specific functions
L.SetGlobal("print", L.NewFunction(func(L *lua.LState) int {
    fmt.Println(L.ToString(1))
    return 0
}))

L.SetGlobal("math_floor", L.NewFunction(func(L *lua.LState) int {
    n := L.ToNumber(1)
    L.Push(lua.LNumber(math.Floor(float64(n))))
    return 1
}))
Oexpose a whole safe subset of a stdlib (like math but not os):
golua.OpenMath(L)   // safe, no system access
lua.OpenString(L) // safe
lua.OpenTable(L)  // safe
// lua.OpenIo(L)  ← don't
// lua.OpenOs(L)  ← don't
// lua.OpenPackage(L) ← don't (blocks require())
Then nuke anything you don't want after loading:
goL.SetGlobal("dofile",    lua.LNil)
L.SetGlobal("loadfile",  lua.LNil)
L.SetGlobal("load",      lua.LNil)
L.SetGlobal("loadstring",lua.LNil)
L.SetGlobal("require",   lua.LNil)










hash the install directory for version tracking 


⚠ python: expected 3.12.0, found 3.12.4
   package may have been modified outside of your tool
   run 'yourpm sync python' to update records

portable ? 
has auto update ? 

track version ? 

can be installed in any dir ? 

needs admin install ? 

can multiple versions be installed side by side ? 


migrate = {
    from = "1.x",
    move = {
        {"/etc/app", "/var/lib/app/config"}
    }
}


    migrate = function()
        -- move old config
    end


scoop packages  →  test suite for the API
if i can express every scoop package in the Lua DSL
i've won
low barrier, high readability




verify final build with a global hash ? 


pause() pause lua code resumed by go cli ?
assert()
print()






PS > Measure-Command {scoop search google}
Results from local buckets...


Days              : 0
Hours             : 0
Minutes           : 0
Seconds           : 26
Milliseconds      : 289
Ticks             : 262892443
TotalDays         : 0.00030427366087963
TotalHours        : 0.00730256786111111
TotalMinutes      : 0.438154071666667
TotalSeconds      : 26.2892443
TotalMilliseconds : 26289.2443


clone/pull git repo     ???
read every single JSON  ???
grep through them all   ???






winget      ->  too rigid, yaml, Microsoft controlled
even more rigid than scoop JSON no logic whatsoever if the installer needs anything custom you're stuck


chocolatey   ->  too complex, .NET dependency, no sandbox
requires .NET runtime
requires NuGet format
PowerShell all the way down
no sandbox whatsoever
complex to contribute to


scoop       ->  sweet spot but PowerShell tax, JSON ceiling



package.hcl
releases/
  1.2.3/
    release.hcl
    script.lua




├───git
│   │   git.hcl
│   │
│   ├───1.2.3
│   │       release.hcl
│   │       script.lua
│   │
│   └───1.3.1
│           release.hcl
│           script.lua
│
└───python
    │   python.hcl
    │
    ├───1.2.3
    │       release.hcl
    │       script.lua
    │
    └───1.3.1
            release.hcl
            script.lua



----------------------------------------------

├───git
│   │   git.hcl
│   │
│   └───versions
│       ├───1.2.3
│       │       release.hcl
│       │       script.lua
│       │
│       └───1.3.1
│               release.hcl
│               script.lua
│
└───python
    │   python.hcl
    │
    └───versions
        ├───1.2.3
        │       release.hcl
        │       script.lua
        │
        └───1.3.1
                release.hcl
                script.lua






-----------------------------------------------------------------------------------------
-- download and verify the installer
installer = fetch "https://github.com/git-for-windows/.../Git-2.44.0-64-bit.exe"
             verify sha256 "c9f1a3..."

-- run it; language knows inno setup flags so you don't have to
run installer as inno:
  components = "gitlfs,assoc"
  editor     = "VIM"
  line_ending = checkout_as_is

-- first-class path manipulation
path.add "%ProgramFiles%\Git\cmd"

-- first-class registry, no cryptic ps syntax
registry.set HKLM\Software\GitForWindows\InstallPath
             = "%ProgramFiles%\Git"

-- assert the install worked before returning
expect shell.output("git --version") contains "2.44.0"
---------------------------------------------------------
-- conditionally choose the right package
let driver_ver = registry.read HKLM\...\NVCUDA\DriverVersion

if driver_ver < "531.00":
  fail "GPU driver too old — update driver first (need 531+)"

-- pick the right build for the machine
let pkg = match system.arch:
  x64   => fetch "https://developer.nvidia.com/.../cuda_12.4_win10_x64.exe"
  arm64 => fail "CUDA not supported on ARM"

run pkg as nsis:
  components = "compiler,libraries,tools"

-- add multiple paths in one go
path.add:
  "%ProgramFiles%\NVIDIA GPU Computing Toolkit\CUDA\v12.4\bin"
  "%ProgramFiles%\NVIDIA GPU Computing Toolkit\CUDA\v12.4\libnvvp"

-- set env vars — native concept, no setx boilerplate
env.set CUDA_PATH = "%ProgramFiles%\NVIDIA GPU Computing Toolkit\CUDA\v12.4"

-- verify a binary exists and runs
expect file.exists "%CUDA_PATH%\bin\nvcc.exe"
expect shell.exit_ok("nvcc --version")
-------------------------------------------------------------

-- filesystem
file.copy    "src\thing.dll"  to="%System32%"
file.extract "app.zip"        to="%ProgramFiles%\MyApp"
file.delete  "%TEMP%\setup.exe"
dir.create   "%AppData%\MyApp\logs"

-- registry
registry.set    HKCU\...\MyApp\Theme = "dark"
registry.read   HKLM\...\Version          -- returns a value
registry.delete HKCU\...\OldKey
registry.exists HKLM\...\MyApp            -- returns bool

-- environment + path
env.set     MY_VAR = "%ProgramFiles%\MyApp"  scope=machine
path.add    "%ProgramFiles%\MyApp\bin"
path.remove "%ProgramFiles%\OldApp\bin"

-- processes + services
service.install "MyApp" exe="%ProgramFiles%\MyApp\svc.exe"
service.start   "MyApp"
service.stop    "MyApp"
shell.run       "netsh advfirewall ..."

-- system inspection
system.arch          -- x64 | arm64
system.windows_ver   -- "11" | "10" | "server2022"
system.is_admin      -- bool
system.reboot_needed -- bool, set automatically after some ops