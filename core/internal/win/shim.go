package win
import (
	_ "embed"

	"runtime"
	"os"
	"fmt"
	"path/filepath"
	"strings"
)

// using the shim from  https://github.com/kiennq/scoop-better-shimexe.git

// files
// assets/shimexe_arm64/shim.exe
// assets/shimexe_x86/shim.exe
// assets/shimexe_x86_64/shim.exe


//go:embed assets/shimexe_x86_64/shim.exe
var shimX64 []byte

//go:embed assets/shimexe_x86/shim.exe
var shimX86 []byte

//go:embed assets/shimexe_arm64/shim.exe
var shimArm64 []byte








func CreateShim(shimPath, target, args string) error {
    err := os.WriteFile(shimPath, shimExe(), 0755)
    if err != nil {
        return fmt.Errorf("creating shim exe: %w", err)
    }

    shimConfig := fmt.Sprintf("path = %s", target)
    if args != "" {
        shimConfig += fmt.Sprintf("\nargs = %s", args)
    }

    configPath := strings.TrimSuffix(shimPath, filepath.Ext(shimPath)) + ".shim"
    err = os.WriteFile(configPath, []byte(shimConfig), 0644)
    if err != nil {
        return fmt.Errorf("creating shim config: %w", err)
    }

    return nil
}


func shimExe() []byte {
    switch runtime.GOARCH {
    case "amd64":
        return shimX64
    case "386":
        return shimX86
    case "arm64":
        return shimArm64
    default:
        return shimX64
    }
}