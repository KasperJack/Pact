package main

import(

	"github.com/kasperjack/pact/core"
	//"os"
	//"path/filepath"
	//"fmt"
	//"strings"
)



type localState struct {
	desktop    string
	cacheDir   string
	installDir string
	repo       string
	lockFile   string
}

func (l *localState) Desktop() string {
	return l.desktop
}

func (l *localState) CacheDir() string {
	return l.cacheDir
}

func (l *localState) InstallDir() string {
	return l.installDir
}

func (l *localState) Repo() string {
	return l.repo
}

func (l *localState) LockFile() string {
	return l.lockFile
}

// repo and lockfile will be use if manager gets a nil Repo or LockFile interface
func NewLocalState(desktop string, cacheDir string, installDir string, repo string, lockFile string) core.LocalState {
	return &localState{
		desktop:    desktop,
		cacheDir:   cacheDir,
		installDir: installDir,
		repo:       repo,
		lockFile:   lockFile,
	}
}


/*
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
*/
