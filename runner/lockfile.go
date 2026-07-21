// lockfile.go — full implementation

package main

import (
    "fmt"
    "os"
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



// GetInstalled returns all installed versions of a package, keyed by version.
// Returns an empty map (not an error) if the package has no installed versions —
// "not installed" is a valid normal state, not a failure

func (f *lockFile) GetInstalled(packageIdentifier string) (map[string]core.LockedPackage) {
    f.mu.RLock()
    defer f.mu.RUnlock()

    result := make(map[string]core.LockedPackage)
    for _, p := range f.content.Packages {
        if p.Name == packageIdentifier {
            result[p.Version] = p
        }
    }
    return result
}

func (f *lockFile) RecordInstall(pkg core.LockedPackage) error {
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

func (f *lockFile) RecordRemove(PackageIdentifier string, version string) error {
    f.mu.Lock()
    defer f.mu.Unlock()

    packages := f.content.Packages
    for i, p := range packages {
        if p.Name == PackageIdentifier {
            f.content.Packages = append(packages[:i], packages[i+1:]...)
            return f.flush()
        }
    }
    return fmt.Errorf("package %q not found", PackageIdentifier)
}

func (f *lockFile) Test() error {
    pkg := core.LockedPackage{
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