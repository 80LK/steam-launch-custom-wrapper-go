package main

import (
	"fmt"
	"os"
	"os/exec"
)

func build() {
	var cmd = exec.Command("go", "build", "-ldflags", "-H=windowsgui", "-o", "slc_wrapper.exe")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error build:", err)
		return
	}
	fmt.Println("Build success.")
}
