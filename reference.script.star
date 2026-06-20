require("dark-archive")
require("ps-extended")



def install(ctx):
    versions = ctx.installed.versions("python")
    
    # if python 3.11 exists keep it, just add this version alongside
    # don't touch its PATH entry
    if "3.11.0" not in versions:
        ctx.env.add_path(ctx.dir())