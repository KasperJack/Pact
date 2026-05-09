package client

import (
	//"path/filepath"
	//"os"
	"strings"
	"url"
	"fmt"
	"io"
	"net/http"

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

	// check if valid gh repo ? 
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
    rawURL, err := gs.buildRawURL(path)
    if err != nil {
        return nil, err
    }

    resp, err := http.Get(rawURL)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch %s: %w", rawURL, err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, rawURL)
    }

    return io.ReadAll(resp.Body)
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