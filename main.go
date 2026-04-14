package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

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
const ERR_REQUIRES_ELEVATION = "The requested operation requires elevation."

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
	if appId == DEFAULT_APP_ID {
		return nil
	}

	fmt.Println("Try init steam sdk")
	if err := os.Setenv(STEAM_APP_ID_ENV_KEY, strconv.Itoa(appId)); err != nil {
		return err
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
	cmd.Stdin = os.Stdin

	fmt.Printf("Try run \"%s\" in directory \"%s\" with arguments [%s]\n", execute, workdir, strings.Join(args, ", "))
	if err := cmd.Start(); err != nil {
		if IsRequiresElevationError(err) {
			fmt.Printf("Run requires elevation, retrying with ElevateRun for \"%s\"\n", execute)
			elevated_cmd := ElevateRun(execute, workdir, args)
			if err := elevated_cmd.Start(); err != nil {
				return err
			}
		}

		return err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	defer signal.Stop(signals)

	select {
	case err := <-done:
		return err
	case <-signals:
		if err := cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			return err
		}

		err := <-done
		if err != nil {
			var exitErr *exec.ExitError
			if !errors.As(err, &exitErr) {
				return err
			}
		}

		return nil
	}
}

func IsRequiresElevationError(err error) bool {
	return err != nil && strings.Contains(err.Error(), ERR_REQUIRES_ELEVATION)
}

func buildPSArray(args []string) string {
	if len(args) == 0 {
		return "NULL"
	}
	escaped := make([]string, len(args))

	for i, a := range args {
		a = strings.ReplaceAll(a, "'", "''")
		escaped[i] = "'" + a + "'"
	}

	return "@(" + strings.Join(escaped, ",") + ")"
}

const CREATE_NO_WINDOW = 0x08000000

func ElevateRun(exe string, workdir string, args []string) *exec.Cmd {
	command := "Start-Process -FilePath '" + exe + "' " +
		"-ArgumentList " + buildPSArray(args) + " " +
		"-WorkingDirectory '" + workdir + "' " +
		`-Verb RunAs`

	cmd := exec.Command(
		"powershell",
		"-Command",
		command,
	)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | CREATE_NO_WINDOW,
	}
	return cmd
}

func main() {
	rootCmd.Flags().IntP(FLAG_APP_ID, P_FLAG_APP_ID, DEFAULT_APP_ID, "usage steam app")
	rootCmd.Flags().StringP(FLAG_WORKDIR, P_FLAG_WORKDIR, DEFAULT_WORKDIR, "usage workdir")
	rootCmd.Execute()
}
