package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ginGen",
	Short: "ginGen is a CLI tool for scaffolding Gin projects",
	Long:  `A fast and flexible CLI tool to generate Gin web projects and add common features.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
