package manager

import (
	"fmt"


	"github.com/kasperjack/pact/core/platform"
	"github.com/kasperjack/pact/core"
	//"runtime"

)




func resolveSource(targetArch string, source core.ReleaseSourceBlock) (core.ReleaseSource, string, error) {
	// targetArch: already verifed to be in "","arm64","x86","x64"

	target := platform.Arch(targetArch)

	//platform.HostArch()

	if target == "" { 
		// fmt.Println("no arch cli override was passed")



		if source.NoArch != nil {
			return *source.NoArch,"noarch",nil 
		}


		if source.Universal != nil {
			return *source.Universal,string(platform.HostArch()),nil 
		}


		switch platform.HostArch() {

			case platform.X64:
				if  source.X64 != nil {
					return *source.X64,string(platform.X64),nil
				}

				if  source.X86 != nil {
					return *source.X86,string(platform.X86),nil
				}
				return core.ReleaseSource{},"",fmt.Errorf("no source found for x64/x86")

			case platform.X86:
				if  source.X86 != nil {
					return *source.X86,string(platform.X86),nil
				}
				return core.ReleaseSource{},"",fmt.Errorf("no source found for x86")
			
			case platform.ARM64:
				if source.ARM64 != nil {
					return *source.ARM64,string(platform.ARM64),nil
				}

				return  core.ReleaseSource{},"",fmt.Errorf("no source found for arm64")

		}

	}else{

			
		//fmt.Println("arch cli override was passed")


		if source.NoArch != nil {
			// warning: arch override has no effect
			return *source.NoArch,"noarch",nil 
		}


		if source.Universal != nil {
			return *source.Universal,targetArch,nil 
		}




		switch platform.HostArch() {

			case platform.X64:
				//fmt.Println(target)
				if target == platform.X64 {
					if  source.X64 != nil {
						return *source.X64,string(platform.X64),nil
					}
					return core.ReleaseSource{},"",fmt.Errorf("no source found for x64")	
				}


				if target == platform.X86 {
					if  source.X86 != nil {
						return *source.X86,string(platform.X86),nil
					}
					return core.ReleaseSource{},"",fmt.Errorf("no source found for x86")	
				}



			case platform.X86:
				if  source.X86 != nil {
					return *source.X86,string(platform.X86),nil
				}
				return core.ReleaseSource{},"",fmt.Errorf("no source found for x86")
			
			case platform.ARM64:
				if source.ARM64 != nil {
					return *source.ARM64,string(platform.ARM64),nil
				}

				return  core.ReleaseSource{},"",fmt.Errorf("no source found for arm64")

		}


	
	}

	return core.ReleaseSource{},"",fmt.Errorf("unexpected error happend")
}