
package main
 
import (
	"fmt"
	"strings"
	"github.com/alecthomas/kong"
	"github.com/kasperjack/pact/core"
	"os"
)




/*
 * TODO:  KNOWN ISSUE: Single-Dash Long Flag Parsing (-arch vs --arch)
 * -------------------------------------------------------------------
 * Behavior:
 * Running `install foo -arch` currently sets:
 *   - Arch = "rch"
 */














 
// command tree add more commands here 
type CLI struct {
	Install InstallCmd `cmd:"" help:"Install a package."`
	List    ListCmd    `cmd:"" help:"List packages."`
}
 
type InstallCmd struct {
	Pkg  string `arg:"" name:"pkg" help:"Package name, optionally pkg@version (e.g. ripgrep@14.1.0)."`
	Arch string `help:"Target architecture (e.g. amd64, arm64)." short:"a"`
	// add more flags here as needed, e.g.:
	// Force  bool `help:"Reinstall even if already installed."`
	// DryRun bool `help:"Show what would happen without doing it." name:"dry-run"`
}
 
func (c *InstallCmd) Run() error {
	pkg, version := parsePkgSpec(c.Pkg)
 
	fmt.Printf("Installing package: %s\n", pkg)
	if version != "" {
		fmt.Printf("  version: %s\n", version)
	} else {
		fmt.Printf("  version: latest\n")
	}
	if c.Arch != "" {
		fmt.Printf("  arch:    %s\n", c.Arch)
	}else {
		fmt.Printf("  arch: auto\n")
	}
 
	linstallDeck(pkg,version,c.Arch)
	return nil
}
 
type ListCmd struct {
	All bool `help:"List all available packages, not just installed ones."`
}
 
func (c *ListCmd) Run() error {
	if c.All {
		fmt.Println("Listing all available packages...")
	} else {
		fmt.Println("Listing installed packages...")
	}
 
	// TODO: actual listing logic goes here
	return nil
}
 
func main() {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("Pcat"),
		kong.Description("a package manager CLI"),
		kong.UsageOnError(),
	)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
 
// parsePkgSpec splits "pkg@version" into its parts.
// "ripgrep" -> ("ripgrep", "")
// "ripgrep@14.1.0" -> ("ripgrep", "14.1.0")
func parsePkgSpec(spec string) (pkg, version string) {
	if i := strings.Index(spec, "@"); i != -1 {
		return spec[:i], spec[i+1:]
	}
	return spec, ""
}








































func linstallDeck(pkg string, version string, arch string){

	
	carch, err := ParseArchFlag(arch)
	if err != nil {
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}


	if err := ValidateArchForHost(carch, core.HostArch()); err != nil {
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}
		


	//fmt.Printf("%s\n",arch)
	
	err = install(pkg,version,carch)

	if err != nil {
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}

	//a := platform.Arch("stringFromCli")
	fmt.Println("everything run ok")
	
}








func ParseArchFlag(s string) (core.Arch, error) {
    switch strings.ToLower(strings.TrimSpace(s)) {
    case "":
        return core.ArchUndefined, nil
    case "x86", "32", "32bit", "386":
        return core.ArchX86, nil
    case "x64", "64", "64bit", "amd64":
        return core.ArchX64, nil
    case "arm64", "aarch64":
        return core.ArchArm64, nil
    default:
        return core.ArchUndefined, fmt.Errorf("unrecognized architecture %q", s)
    }
}





func ValidateArchForHost(requested, host core.Arch) error {
    if requested == core.ArchUndefined {
        return nil 
    }
	
    if requested == host {
        return nil
    }
   
    if requested == core.ArchX86 && host == core.ArchX64{
        return nil
    }

    return fmt.Errorf("architecture %s is not supported on this host (%s) — try omitting --arch to auto-select, or use a compatible architecture",requested.String(),host.String())
}







/*
func validateArch(targetArch platform.Arch) error {
    
	host, err := platform.HostArch()

	if err != nil {
		return err 
	}
	
    switch host {

		case platform.ARM64:
			if targetArch != platform.ARM64 {
				return fmt.Errorf("cannot install %s binary on an arm64 host", targetArch)
			}
	


		case platform.X64:
			if targetArch != platform.X64 && targetArch != platform.X86 {
				return fmt.Errorf("cannot install %s binary on an x64 host", targetArch)
			}
	

		case platform.X86:
			if targetArch != platform.X86 {
				return fmt.Errorf("cannot install %s binary on an x86 host", targetArch)
			}
		default:
    		return fmt.Errorf("unsupported host architecture: %s", host)
	}

	return nil
}*/








