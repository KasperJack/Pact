package client


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

	b,err := rc.source.Fetch(path)
		if err != nil {panic(err)}

	print(string(b))
}