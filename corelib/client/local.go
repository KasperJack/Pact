package client


import (
	"path/filepath"
	"os"

)

func NewLFilesystemSource(rootdir string) *FilesystemSource {
    return &FilesystemSource{RootDir: rootdir}
}


type FilesystemSource struct {
	RootDir string
}


func (fss* FilesystemSource) Fetch (path string) ([]byte, error) {

	
	data, err := os.ReadFile(filepath.Join(fss.RootDir,path))
    if err != nil {
        return nil, err
    }

    return data, nil
}
