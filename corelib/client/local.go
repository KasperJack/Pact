package client


import (
	"path/filepath"
	"os"
    "fmt"
    "errors"
    "io/fs"

)
var ErrNotFound = errors.New("file not found")


func NewLFilesystemSource(rootdir string) *FilesystemSource {
    // paic if root dir is not there ? 
    return &FilesystemSource{RootDir: rootdir}
}


type FilesystemSource struct {
	RootDir string
}


func (fss* FilesystemSource) Fetch (path string) ([]byte, error) {

	
    data, err := os.ReadFile(filepath.Join(fss.RootDir, path))
    if err != nil {
        if errors.Is(err, fs.ErrNotExist) {
            return nil, fmt.Errorf("%w: %s", ErrNotFound, path)
        }
        return nil, err
    }

    return data, nil
}
