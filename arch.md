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