package manager

import (
    "archive/zip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
)

// extract an archive in a folder and delete the archive
func Extract(dir string) error { //NM: tar, 7zip needs and exteranl lib for managing archives 
    
    entries, err := os.ReadDir(dir)
    if err != nil {
        return err
    }

    var zipPath string
    for _, e := range entries {
        if filepath.Ext(e.Name()) == ".zip" {
            zipPath = filepath.Join(dir, e.Name())
            break
        }
    }

    if zipPath == "" {
        return fmt.Errorf("no zip file found in %s", dir)
    }

    // extract
    r, err := zip.OpenReader(zipPath)
    if err != nil {
        return err
    }
    defer r.Close()

    for _, f := range r.File {
        dest := filepath.Join(dir, f.Name)

        if !strings.HasPrefix(dest, filepath.Clean(dir)+string(os.PathSeparator)) {
            return fmt.Errorf("illegal path in zip: %s", f.Name)
        }

        if f.FileInfo().IsDir() {
            os.MkdirAll(dest, f.Mode())
            continue
        }

        if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
            return err
        }

        out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
        if err != nil {
            return err
        }

        rc, err := f.Open()
        if err != nil {
            out.Close()
            return err
        }

        _, err = io.Copy(out, rc)
        out.Close()
        rc.Close()
        if err != nil {
            return err
        }
    }

    // close before delete
    r.Close()

    // delete the zip
    return os.Remove(zipPath)
}