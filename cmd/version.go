package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const CLIVersion = "0.1.0-mvp"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ginGen",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ginGen CLI Tool v%s\n", CLIVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
