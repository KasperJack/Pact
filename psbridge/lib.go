package psbridge

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ASS() {
	// get absolute paths so pwsh can find everything
	sessionScript, _ := filepath.Abs("session.ps1")
	userScript, _    := filepath.Abs("script.ps1")
    // "--memory", "256m",
	cmd := exec.Command("powershell",
		"-NonInteractive",
		"-NoProfile",
		"-File", sessionScript,
		"-ScriptPath", userScript,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("running script...")
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("done")
}

func GG(){
	fmt.Println("gg")
}