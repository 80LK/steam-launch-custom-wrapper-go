package commands

import "github.com/spf13/cobra"

var RootCommand = &cobra.Command{
	Use:   "devtool",
	Short: "Devtools for SLC-Wrapper",
	Long:  "Devtools for SLC-Wrapper",
}
