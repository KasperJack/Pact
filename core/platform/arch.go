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

func HostArch() Arch {
    switch runtime.GOARCH {
    case "386":
        return X86
    case "amd64":
        return X64
    default: // arm64
        return ARM64
    }
}


// initail setp 
func ParseArch(s string) (Arch, error) {
    switch Arch(s) {
    case X86, X64, ARM64:
        return Arch(s), nil
    default:
        return "", fmt.Errorf("unknown arch %q, must be one of: x86, x64, arm64", s)
    }
}




func (target Arch) IsCompatibleWith(host Arch) bool { 
        switch host {

            case ARM64:
                if target != ARM64 {
                    return false
                }

            case X64:
                if target!= X64 && target != X86 {
                    return false
                }
        

            case X86:
                if target != X86 {
                    return false
                }
            default:
                return false

	    }
        return true
}








func (target Arch) ValidateForHost() error {
    host := HostArch()
    
    if !target.IsCompatibleWith(host) {
        return fmt.Errorf("cannot install %s binary on an %s host", target, host)
    }

    return nil
}