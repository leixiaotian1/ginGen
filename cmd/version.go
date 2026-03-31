package cmd

import (
	"fmt"

	"github.com/leixiaotian1/ginGen/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ginGen",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ginGen %s (commit %s, built %s)\n", version.Version, version.Commit, version.BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
