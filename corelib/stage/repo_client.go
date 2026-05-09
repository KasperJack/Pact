package stage


type RepositorySource interface {
	Fetch (path string) ([]byte, error) 

}