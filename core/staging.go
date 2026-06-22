package core
import (
	"os"
	"path/filepath"
)
// core/staging.go
type StagingArea struct {
    root string // C:\Users\Aya\Desktop\pact\bin\staging
}

func NewStagingArea(root string) *StagingArea {
    return &StagingArea{root: root}
}

func (s *StagingArea) Dir(pkg, version string) string {
    return filepath.Join(s.root, pkg, version)
}

func (s *StagingArea) Prepare(pkg, version string) (string, error) {
    dir := s.Dir(pkg, version)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return "", err
    }
    return dir, nil
}

func (s *StagingArea) Clear() error {
    if err := os.RemoveAll(s.root); err != nil {
        return err
    }

    return os.MkdirAll(s.root, 0755) //RE: refactor a propper error 
}