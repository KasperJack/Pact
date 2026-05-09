package client


type RepositorySource interface {
	Fetch (path string) ([]byte, error) 
}

type RepoClient struct {
	source RepositorySource
}


func NewRepoClient(source RepositorySource) *RepoClient {
    return &RepoClient{source: source}
}