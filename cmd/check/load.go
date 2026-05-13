package main

import (
    "os"
    "path/filepath"
)


func loadFile(path string) ([]byte, error) {
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