package install

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "net/http"
    "os"
    "path"
    "path/filepath"
)

// chnage to private func 
func download(source string, hash string, targetDir string) error { // download check hash delete file hash missmatch 
    resp, err := http.Get(source)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("download failed: %s", resp.Status)
    }

    filename := path.Base(source)  // should get Content-Disposition 
    dest := filepath.Join(targetDir, filename)

    f, err := os.Create(dest)
    if err != nil {
        return err
    }
    defer f.Close()

    h := sha256.New()
	fmt.Println("downloading file")
    _, err = io.Copy(io.MultiWriter(f, h), resp.Body)
    if err != nil {
        return err
    }

    actual := hex.EncodeToString(h.Sum(nil))
    if actual != hash {
        //os.Remove(dest)
        return fmt.Errorf("checksum mismatch: expected %s got %s", hash, actual)
    }

    return nil
}