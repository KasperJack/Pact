package client

import (
    "encoding/base64"
    "encoding/json"
    "net/http"
    "strings"
    "fmt"
)


type githubContent struct {
    Content  string `json:"content"`
    Encoding string `json:"encoding"`
}


type GithubAPISource struct {
    User   string
    Repo   string
    Branch string
    Token  string
}





func (gs *GithubAPISource) Fetch(path string) ([]byte, error) {
    url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s",
        gs.User, gs.Repo, strings.TrimLeft(path, "/"), gs.Branch,
    )

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "Bearer "+gs.Token)
    req.Header.Set("Accept", "application/vnd.github+json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
    }

    var content githubContent
    if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    decoded, err := base64.StdEncoding.DecodeString(
        strings.ReplaceAll(content.Content, "\n", ""),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to decode content: %w", err)
    }

    return decoded, nil
}