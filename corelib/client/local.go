package client


import (
	"path/filepath"
	"os"

)


type FilesystemSource struct {
	RootDir string
}
// methods here 

func (fss* FilesystemSource) Fetch (path string) ([]byte, error) {

	
	data, err := os.ReadFile(filepath.Join(fss.RootDir,path))
    if err != nil {
        return nil, err
    }

    return data, nil
}
