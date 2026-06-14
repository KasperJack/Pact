package platform
import (
	"runtime"
	"fmt"

)


type Arch string

const (
    X86   Arch = "x86"
    X64   Arch = "x64"
    ARM64 Arch = "arm64"
)

func HostArch() (Arch, error) {
    switch runtime.GOARCH {
    case "386":
        return X86, nil

    case "amd64":
        return X64, nil

    case "arm64":
        return ARM64, nil

    default:
        return "", fmt.Errorf(
            "unsupported architecture: %s",
            runtime.GOARCH,
        )
    }
}
