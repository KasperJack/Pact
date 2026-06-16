// lockfile.go — full implementation

package main

import (
    "fmt"
    "os"
    "sync"
	"time"

    "github.com/kasperjack/pact/core"
    "github.com/kasperjack/pact/core/model"
    "github.com/kasperjack/pact/core/parce"
)

func NewLockFile(filePath string) (core.LockFile, error) {
    f, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("error loading the lock file") //RF:E
    }

    c, err := parce.LockFile(f)
    if err != nil {
        return nil, err
    }

    return &lockFile{path: filePath, content: c}, nil
}

type lockFile struct {
    mu      sync.RWMutex
    path    string
    content model.LockFile
}

func (f *lockFile) GetInstalled(pkg string) (model.LockedPackage, error) {
    f.mu.RLock()
    defer f.mu.RUnlock()

    for _, p := range f.content.Packages {
        if p.Name == pkg {
            return p, nil
        }
    }
    return model.LockedPackage{}, fmt.Errorf("package %q is not installed", pkg)
}

func (f *lockFile) RecordInstall(pkg model.LockedPackage) error {
    f.mu.Lock()
    defer f.mu.Unlock()

    // replace if already present
    for i, p := range f.content.Packages {
        if p.Name == pkg.Name {
            f.content.Packages[i] = pkg
            return f.flush()
        }
    }
    f.content.Packages = append(f.content.Packages, pkg)
    return f.flush()
}

func (f *lockFile) RecordRemove(pkg string) error {
    f.mu.Lock()
    defer f.mu.Unlock()

    packages := f.content.Packages
    for i, p := range packages {
        if p.Name == pkg {
            f.content.Packages = append(packages[:i], packages[i+1:]...)
            return f.flush()
        }
    }
    return fmt.Errorf("package %q not found", pkg)
}

func (f *lockFile) Test() error {
    pkg := model.LockedPackage{
        Name:        "test-package2",
        Version:     "1.0.0",
        InstalledAt: time.Now().Format(time.RFC3339),
        InstallDir:  "/usr/local/test-package",
    }
    return f.RecordInstall(pkg)
}


func (f *lockFile) flush() error {
    data := parce.WriteLockFile(f.content)
    return os.WriteFile(f.path, data, 0644)
}