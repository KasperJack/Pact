package client
import (
    "errors"
	"path/filepath"
)

type RepositorySource interface {
	Fetch (path string) ([]byte, error) 
}


func NewRepoClient(source RepositorySource) *RepoClient {
    return &RepoClient{source: source}
}

type RepoClient struct {
	source RepositorySource
}

func (rc *RepoClient) PackageExists(packageName string) bool {

	//  /pa/packageName/package.toml
		path := filepath.Join(packageName[:2], packageName, "package.toml")

	_,err := rc.source.Fetch(path)
		if err != nil {

			if errors.Is(err,ErrNotFound) {
				return false
			}
			panic(err)	
		}

	
	return true
}