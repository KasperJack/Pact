manifest.lua  →  gopher-lua  →  Go functions  →  OS
   (portable)      (bridge)      (API)     (specific ps/sh)


//
verify = function(ctx)
    ctx.assert_file(ctx.dir .. "/python.exe")
    ctx.assert_shim("python")
    ctx.assert_runs("python --version")
end

or maybe do assersion from the go side 