package model


type ReleaseSource struct {
    URL    string `hcl:"url"`
    SHA256 string `hcl:"sha256"`
}

type ReleaseSourceBlock struct {
    X64   *ReleaseSource `hcl:"x64,block"`
    ARM64 *ReleaseSource `hcl:"arm64,block"`
    X86   *ReleaseSource `hcl:"x86,block"`
}

type Release struct {
    Version string             `hcl:"version"`
    Source  ReleaseSourceBlock `hcl:"source,block"`
}