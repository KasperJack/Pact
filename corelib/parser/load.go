package parser

import (
    "os"
    "path/filepath"
)

// LoadFile reads the contents of the file at the given path and returns it as a byte slice.
// The path is converted to an absolute path before reading.
//
// It returns an error if the path cannot be resolved or if the file cannot be read.
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