package main

import(

	"github.com/kasperjack/pact/core"
	"os"
	"path/filepath"


)

func NewLocalState() *core.LocalState {

	
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exeDir := filepath.Dir(exePath)


	userProfile := os.Getenv("USERPROFILE")
	appData := os.Getenv("APPDATA")
	localAppData := os.Getenv("LOCALAPPDATA")



	
	publicProfile := os.Getenv("PUBLIC")
	programFiles := os.Getenv("PROGRAMFILES")
	programData := os.Getenv("PROGRAMDATA")





	return &core.LocalState{
		CacheDir: filepath.Join(exeDir, "cache"),
		Repo:     filepath.Join(exeDir, "repo"),

		UserLockFile:   filepath.Join(localAppData, "pact", "user.hcl"),
		SystemLockFile: filepath.Join(programFiles, "pact", "system.hcl"),

		UserDesktop:   filepath.Join(userProfile, "Desktop"),
		PublicDesktop: filepath.Join(publicProfile, "Desktop"),

		UserPackagesDir:   filepath.Join(localAppData, "pact", "packages"),
		SystemPackagesDir: filepath.Join(programFiles, "pact", "packages"),


		UserProfile: userProfile,
		AppData: appData,
		LocalAppData: localAppData,


		PublicProfile: publicProfile,
		ProgramFiles: programFiles,
		ProgramData: programData,
	}
}




