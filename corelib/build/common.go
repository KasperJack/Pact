package build

import (
    "os"
    "github.com/BurntSushi/toml"
	"path/filepath"
)



func writeToml(path string, v any) error {
    if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
        return err
    }

    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()

    return toml.NewEncoder(f).Encode(v)
}
