package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const WD_FLAG = "--wd"
const WD_FLAG_EQ = WD_FLAG + "="

// go:embed go.ver
var VERSION string

func prepareArgs() (bool, string, string, []string) {
	var os_args = os.Args[1:]

	var execute string = ""
	var workdir string = ""
	var args = []string{}

	if len(os_args) == 0 {
		return false, execute, workdir, args
	}

	var nextWorkDir = false
	var foundedExecute = false
	for _, arg := range os_args {
		if arg == WD_FLAG {
			nextWorkDir = true
			continue
		}

		if nextWorkDir {
			nextWorkDir = false
			workdir, _ = filepath.Abs(arg)
			continue
		}

		if strings.HasPrefix(arg, WD_FLAG_EQ) {
			arg, _ = strings.CutPrefix(arg, WD_FLAG_EQ)
			workdir, _ = filepath.Abs(arg)
			continue
		}

		if !foundedExecute {
			execute, _ = filepath.Abs(arg)
			if workdir == "" {
				workdir = filepath.Dir(execute)
			}
			foundedExecute = true
			continue
		}

		args = append(args, arg)
	}

	return foundedExecute, execute, workdir, args
}

func main() {
	var parsed, execute, workdir, args = prepareArgs()
	if !parsed {
		fmt.Printf("Steam-Launch-Custom-Wrapper %s\nUsage:\n\twrapper [execute] <--wd=[wordir]> <launch args>", VERSION)
		return
	}

	var cmd = exec.Command(execute, args...)
	cmd.Dir = workdir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf(`Run "%s" in directory "%s" with arguments [%s]`, execute, workdir, strings.Join(args, ", "))
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
