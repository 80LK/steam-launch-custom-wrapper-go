package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hajimehoshi/go-steamworks"
	"github.com/spf13/cobra"
)

//go:embed go.ver
var VERSION string

const STEAM_APP_ID_ENV_KEY = "SteamAppId"
const FLAG_WORKDIR = "workdir"
const P_FLAG_WORKDIR = "w"
const DEFAULT_WORKDIR = ""
const FLAG_APP_ID = "app"
const P_FLAG_APP_ID = "a"
const DEFAULT_APP_ID = 0

var rootCmd = &cobra.Command{
	Use:   "slc_wrapper <executable> [arguments...]",
	Short: "slc_wrapper - Steam-Launch-Custom Wrapper",
	Long: `slc_wrapper	Steam-Launch-Custom Wrapper.
	Usage for startup apps how steam game`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		AppId, err := cmd.Flags().GetInt(FLAG_APP_ID)
		if err != nil {
			return err
		}

		WorkDir, err := cmd.Flags().GetString(FLAG_WORKDIR)
		if err != nil {
			return err
		}

		Executable, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		if WorkDir == DEFAULT_WORKDIR {
			WorkDir = filepath.Dir(Executable)
		}
		Args := args[1:]

		if err := ValidateExecute(Executable); err != nil {
			return err
		}
		if err := ValidateWorkDir(WorkDir); err != nil {
			return err
		}

		if err := SteamWorksInit(AppId); err != nil {
			return err
		}

		if err := Run(Executable, WorkDir, Args); err != nil {
			return err
		}

		return nil
	},
}

func CheckExsist(_t string, path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%s path \"%s\" is not exsist", _t, path)
		}
		return nil, err
	}

	return info, nil
}

func ValidateExecute(path string) error {
	const _t = "Executable"
	info, err := CheckExsist(_t, path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("%s path \"%s\" is Directory", _t, path)
	}

	return nil
}

func ValidateWorkDir(path string) error {
	const _t = "WorkDir"
	info, err := CheckExsist(_t, path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	return fmt.Errorf("%s path \"%s\" is not Directory", _t, path)
}

func SteamWorksInit(appId int) error {
	fmt.Println("Try init steam sdk")
	if appId != DEFAULT_APP_ID {
		if err := os.Setenv(STEAM_APP_ID_ENV_KEY, strconv.Itoa(appId)); err != nil {
			return err
		}
	}

	if err := steamworks.Init(); err != nil {
		return err
	}

	return nil
}

func Run(execute string, workdir string, args []string) error {
	var cmd = exec.Command(execute, args...)
	cmd.Dir = workdir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Try run \"%s\" in directory \"%s\" with arguments [%s]\n", execute, workdir, strings.Join(args, ", "))
	return cmd.Run()
}

func main() {
	rootCmd.Flags().IntP(FLAG_APP_ID, P_FLAG_APP_ID, DEFAULT_APP_ID, "usage steam app")
	rootCmd.Flags().StringP(FLAG_WORKDIR, P_FLAG_WORKDIR, DEFAULT_WORKDIR, "usage workdir")
	rootCmd.Execute()
}
