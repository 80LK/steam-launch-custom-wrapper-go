package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

const _OUTPUT_FLAG = "output"
const _OUTPUT_FLAG_P = "o"

var BuildCommand = &cobra.Command{
	Use:   "build",
	Short: "Build project",
	Long:  "Build project",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, err := cmd.Flags().GetString(_OUTPUT_FLAG)
		if err != nil {
			return err
		}

		var build = exec.Command("go", "build", "-o", output)

		build.Stdout = os.Stdout
		build.Stderr = os.Stderr

		if err := build.Run(); err != nil {
			return err
		}

		fmt.Println("Build success.")
		return nil
	},
}

func addBuildFlag(cmd *cobra.Command) {
	cmd.Flags().StringP(_OUTPUT_FLAG, _OUTPUT_FLAG_P, "slc_wrapper.exe", "Output file name")
}

func init() {
	addBuildFlag(BuildCommand)
	RootCommand.AddCommand(BuildCommand)
}
