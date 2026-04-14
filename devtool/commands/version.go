package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"devtool/version"

	"github.com/spf13/cobra"
)

const VERSION_FILE = "go.ver"

func getVersion(version string) (int, int, int, error) {
	parts := strings.Split(version, ".")
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, err
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, err
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, err
	}

	return major, minor, patch, nil
}

func getVersionFromFile() (*version.Version, error) {
	var data, err = os.ReadFile(VERSION_FILE)
	if err != nil {
		return nil, err
	}

	return version.Parse(string(data))
}

var _VERSION_KEYWORDS = []string{"patch", "minor", "major"}

var version_keywordsMap map[string]bool = nil

func validKeyword(value string) bool {
	if version_keywordsMap == nil {
		version_keywordsMap = make(map[string]bool)
		for _, v := range _VERSION_KEYWORDS {
			version_keywordsMap[v] = true
		}
	}

	return version_keywordsMap[value]
}

var VersionCommand = &cobra.Command{
	Use:   "version <" + strings.Join(_VERSION_KEYWORDS, "|") + "|version>",
	Short: "Short version",
	Long:  "Long version",
	Args: func(cmd *cobra.Command, args []string) error {
		raw_version := strings.ToLower(args[0])
		if validKeyword(raw_version) {
			return nil
		}

		_, err := version.Parse(raw_version)
		if err == nil {
			return nil
		}

		return fmt.Errorf("Not valid version format: %s", raw_version)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		gitStatus, err := exec.Command("git", "status", "--porcelain=v1", "-uno").Output()
		if err != nil {
			return err
		}
		if len(gitStatus) != 0 {
			return fmt.Errorf("Git working directory not clean.")
		}

		raw_version := strings.ToLower(args[0])

		if validKeyword(raw_version) {
			current_version, err := getVersionFromFile()
			if err != nil {
				return err
			}

			switch raw_version {
			case "patch":
				current_version.Patch += 1
			case "minor":
				current_version.Patch = 0
				current_version.Minor += 1
			case "major":
				current_version.Patch = 0
				current_version.Minor = 0
				current_version.Major += 1
			}

			raw_version = current_version.String()
		}

		err = os.WriteFile(VERSION_FILE, []byte(raw_version), 0644)
		if err != nil {
			return err
		}

		err = exec.Command("git", "add", VERSION_FILE).Run()
		if err != nil {
			return err
		}
		err = exec.Command("git", "commit", "-m", raw_version).Run()
		if err != nil {
			return err
		}
		err = exec.Command("git", "tag", "v"+raw_version).Run()
		if err != nil {
			return err
		}

		fmt.Println("You need to build the project yourself")
		return nil
	},
}

func init() {
	RootCommand.AddCommand(VersionCommand)
}
