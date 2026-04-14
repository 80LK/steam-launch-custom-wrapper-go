package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const VERSION_FILE = "go.ver"

func getVersion() (int, int, int, error) {
	var data, err = os.ReadFile(VERSION_FILE)
	if err != nil {
		return 0, 0, 0, err
	}

	var version string = string(data)
	var parts = strings.Split(version, ".")
	var major, _ = strconv.Atoi(parts[0])
	var minor, _ = strconv.Atoi(parts[1])
	var patch, _ = strconv.Atoi(parts[2])

	return major, minor, patch, nil
}

type VersionPart string

const (
	Patch VersionPart = "patch"
	Major VersionPart = "Major"
	Minor VersionPart = "Minor"
)

func version(m VersionPart) {
	var gitStatus, _ = exec.Command("git", "status", "--porcelain=v1", "-uno").Output()
	if len(gitStatus) != 0 {
		fmt.Println("Git working directory not clean.")
		return
	}

	var major, minor, patch, _ = getVersion()

	switch m {
	case Patch:
		patch += 1
	case Minor:
		patch = 0
		minor += 1
	case Major:
		patch = 0
		minor = 0
		major += 1
	}
	var version = strconv.Itoa(major) + "." + strconv.Itoa(minor) + "." + strconv.Itoa(patch)
	os.WriteFile(VERSION_FILE, []byte(version), 0644)

	exec.Command("git", "add", VERSION_FILE).Run()
	exec.Command("git", "commit", "-m", version).Run()
	exec.Command("git", "tag", "v"+version).Run()

	build()
}
