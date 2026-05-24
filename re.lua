version     = "3.12.0"


source = {
    x64 = {
        url = "https://github.com/windirstat/windirstat/releases/download/release%2Fv2.6.1/WinDirStat.zip",
        sha256 = "42de19daff2bcfa0cc177626593edbeda411a745ca2079d3f28aec4cffe067cd"
    },
    arm64 = {
        url = "https://github.com/windirstat/windirstat/releases/download/release%2Fv2.6.1/WinDirStat.zip",
        sha256 = "42de19daff2bcfa0cc177626593edbeda411a745ca2079d3f28aec4cffe067cd"
    },
    x86 = {
        url = "https://github.com/windirstat/windirstat/releases/download/release%2Fv2.6.1/WinDirStat.zip",
        sha256 = "42de19daff2bcfa0cc177626593edbeda411a745ca2079d3f28aec4cffe067cd"
    }
}

install = function(ctx)
    ctx.extract(ctx.dist(), ctx.staging())

    if ctx.os.x64() then
        ctx.move(ctx.path.join(ctx.staging(), "x64/*"), ctx.dir())
    elseif ctx.os.x86() then
        ctx.move(ctx.path.join(ctx.staging(), "x86/*"), ctx.dir())
    elseif ctx.os.arm64() then
        ctx.move(ctx.path.join(ctx.staging(), "arm64/*"), ctx.dir())
    end
end

post_install = function(ctx)
    ctx.shortcut(ctx.path.join(ctx.dir(), "WinDirStat.exe"), "WinDirStat")
end