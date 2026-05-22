Package {
    name = "git",
    version = "2.48.1",

    description = "Distributed version control system",
    homepage = "https://git-scm.com",

    sources = {
        windows = {
            url = "https://github.com/git-for-windows/git/releases/download/v2.48.1.windows.1/Git-2.48.1-64-bit.exe",
            sha256 = "abc123..."
        },

        linux = {
            url = "https://mirrors.edge.kernel.org/pub/software/scm/git/git-2.48.1.tar.gz",
            sha256 = "def456..."
        }
    } }