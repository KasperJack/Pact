package stage

import (
    "os"
    "path/filepath"
)

func LoadFile(path string) ([]byte, error) {
    abs, err := filepath.Abs(path)
    if err != nil {
        return nil, err
    }

    data, err := os.ReadFile(abs)
    if err != nil {
        return nil, err
    }

    return data, nil
}