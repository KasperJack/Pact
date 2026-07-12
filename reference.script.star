require("dark-archive")
require("ps-extended")



def install(ctx):
    versions = ctx.installed.versions("python")

    if "3.11.0" not in versions:
        ctx.env.add_path(ctx.dir())