package main

import(

	"github.com/kasperjack/pact/core"
	"os"
	"path/filepath"
	"fmt"
	"strings"
)


type localState struct {
	baseDir string
}



func NewLocalState(baseDir string) core.LocalState {
 return &localState{
	baseDir: baseDir,
 }
}




func (l *localState) PackageExists(PackageIdentifier string) (bool, error) {

    return true, nil
}



func (l *localState) CreateInstallDir(PackageIdentifier string, version string) (string, error) {
	dir := filepath.Join(l.baseDir, PackageIdentifier, version)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create install dir for %s@%s: %w", PackageIdentifier, version, err)
	}

	return dir, nil
}


func (l *localState) RemovePackageDir(path string) error {

	// add safety checks
	if path == "" {
		return fmt.Errorf("cannot remove package dir: path is empty")
	}

	absBase, err := filepath.Abs(l.baseDir)
	if err != nil {
		return fmt.Errorf("failed to resolve base dir: %w", err)
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path %q: %w", path, err)
	}
	rel, err := filepath.Rel(absBase, absPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("refusing to remove path outside base dir: %q", path)
	}

	if err := os.RemoveAll(absPath); err != nil {
		return fmt.Errorf("failed to remove package dir %q: %w", path, err)
	}

	return nil
}

