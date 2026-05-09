package client

import (
	//"path/filepath"
	//"os"
	"strings"
	"url"
	"fmt"

)


type GithubSource struct {
    RepoURL string
	Branch string // should defult to "main" if not passed 
	// what happens if i create a struct without passing all value do they take the 0 value ? 
	User string
	Repo string

}

func NewGithubSource(repoURL string, branch string) (*GithubSource, error) {
    parsed, err := url.Parse(strings.TrimRight(repoURL, "/"))
    if err != nil {
        return nil, fmt.Errorf("invalid url: %s", repoURL)
    }

    parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
    if len(parts) < 2 {
        return nil, fmt.Errorf("cannot parse user/repo from URL: %s", repoURL)
    }

    if branch == "" {
        branch = "main"
    }

    return &GithubSource{
        RepoURL: repoURL,
        User:    parts[0],
        Repo:    parts[1],
        Branch:  branch,
    }, nil
}


func (gs *GithubSource) Fetch(path string) ([]byte, error) {
		_, rawUrl := gs.buildRawURL(path)
}


func (gs *GithubSource) buildRawURL(path string) (string, error) {
    if !strings.HasPrefix(path, "/") {
        return "", fmt.Errorf("path must be absolute: %s", path)
    }

    clean := strings.TrimLeft(path, "/")

    return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s",
        gs.User,
        gs.Repo,
        gs.Branch,
        clean,
    ), nil
}