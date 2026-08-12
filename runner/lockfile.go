// lockfile.go — full implementation

package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
	"time"

    "github.com/kasperjack/pact/core"
    "github.com/kasperjack/pact/core/parce"
)
                                    //interface                
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
    content core.LockFileC
}



func (f *lockFile) GetInstalled(packageIdentifier string) (core.LockedPackage, bool) {
    f.mu.RLock()
    defer f.mu.RUnlock()

    for _, p := range f.content.Packages {
        if p.Identifier == packageIdentifier {
            return p, true
        }
    }
    return core.LockedPackage{}, false
}


func (f *lockFile) RecordInstall(pkg core.LockedPackage) error {
    f.mu.Lock()
    defer f.mu.Unlock()

    // replace if already present
    for i, p := range f.content.Packages {
        if p.Identifier == pkg.Identifier {
            f.content.Packages[i] = pkg
            return f.flush()
        }
    }
    f.content.Packages = append(f.content.Packages, pkg)
    return f.flush()
}


func (f *lockFile) RecordRemove(packageIdentifier string) error {
    f.mu.Lock()
    defer f.mu.Unlock()

    packages := f.content.Packages
    for i, p := range packages {
        if p.Identifier == packageIdentifier {
            f.content.Packages = append(packages[:i], packages[i+1:]...)
            return f.flush()
        }
    }
    return fmt.Errorf("%w: %s",core.ErrPkgNotFound,packageIdentifier)
}



func (f *lockFile) Test() error {

    testPath := filepath.Join("ass,hole","test-install-dir")

    pkg := core.LockedPackage{
        Identifier:        "test-package",
        Version:     "1.0.0",
        InstalledAt: time.Now().Format(time.RFC3339),
        InstallDir:  testPath,
    }
    return f.RecordInstall(pkg)
}


func (f *lockFile) flush() error {
    data := parce.WriteLockFile(f.content)
    return os.WriteFile(f.path, data, 0644)
}