package model

type LockedPackage struct {
    Name        string `hcl:"name,label"`
    Version     string `hcl:"version"`
    InstalledAt string `hcl:"installed_at"`
    InstallDir  string `hcl:"install_dir"`
}

type LockFile struct {
    Packages []LockedPackage `hcl:"package,block"`
}