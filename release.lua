python {
    
    version     = "3.12.0",

    source = {
        url    = "https://www.python.org/ftp/python/3.12.0/python-3.12.0-amd64.exe",
        sha256 = "abc123...",
    },

    pre_install = function()
        local arch_label, bit = "64-bit", "64"

        if arch() == "32bit" then
            arch_label, bit = "32-bit", "32"
        elseif arch() == "arm64" then
            arch_label = "ARM64"
        end

        local py_root    = dir():gsub("\\", "\\\\")
        local py_version = version():match("^(%d+%.%d+)")

     

        for _, reg_file in ipairs({ "install-pep-514.reg", "uninstall-pep-514.reg" }) do
            replace_in_file(script(reg_file), dir(reg_file), {
                ["$py_root"]         = py_root,
                ["$py_version"]      = py_version,
                ["$py_fullversion"]  = version(),
                ["$py_cleanVersion"] = version():gsub("%.", ""),
                ["$py_archLabel"]    = arch_label,
                ["$py_arch"]         = bit,
                ["HKEY_CURRENT_USER"] = is_global() and "HKEY_LOCAL_MACHINE" or "HKEY_CURRENT_USER",
            })
        end
    end,

    install = function()
        expand_archive(dir("setup.exe"), dir("_tmp"))

        for _, f in ipairs({ "path.msi", "pip.msi" }) do
            remove(dir("_tmp/AttachedContainer/" .. f))
        end

        for _, msi in ipairs(glob(dir("_tmp/AttachedContainer/*.msi"))) do
            if basename(msi) ~= "appendpath" then
                expand_msi(msi, dir())
            end
        end

        remove(dir("_tmp"), dir("setup.exe"))

        if is_global() then
            local pathext = env_get("PATHEXT"):gsub(";%.PYW?", "")
            env_set("PATHEXT", pathext .. ";.PY;.PYW")
        end
    end,

    post_install = function()
        sh("python -E -s -m ensurepip -U --default-pip")
    end,

    uninstall = function()
        if is_global() then
            local pathext = env_get("PATHEXT"):gsub(";%.PYW?", "")
            env_set("PATHEXT", pathext)
        end
    end,
}






install = function()
        printFrlomGo("fd")
  
    end,