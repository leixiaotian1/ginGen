package cmd

import (
	"fmt"
	"ginGen/internal/generator"
	"ginGen/internal/utils"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var modulePath string

var newCmd = &cobra.Command{
	Use:   "new <project_name>",
	Short: "Create a new Gin project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		targetModulePath := modulePath

		if targetModulePath == "" {
			// Default module path, e.g., "project_name" or prompt user for github.com/user/project
			// For MVP, we'll use projectName directly. Consider making this more robust later.
			targetModulePath = projectName
			fmt.Printf("Module path not specified, defaulting to: %s\n", targetModulePath)
			fmt.Println("You can specify a custom module path with --module, e.g., github.com/youruser/yourproject")
		}

		fmt.Printf("Creating new Gin project: %s (Module: %s)\n", projectName, targetModulePath)

		// 1. Create project directory
		if err := os.MkdirAll(projectName, 0755); err != nil {
			log.Fatalf("Error creating project directory %s: %v", projectName, err)
		}
		fmt.Printf("Created project directory: %s\n", projectName)

		// Use absolute path for projectPath for subsequent operations
		absProjectPath, err := filepath.Abs(projectName)
		if err != nil {
			log.Fatalf("Error getting absolute path for %s: %v", projectName, err)
		}

		// 2. Generate project structure and files from templates
		templateData := generator.TemplateData{
			ProjectName: projectName,
			ModulePath:  targetModulePath,
		}
		if err := generator.GenerateProjectStructure(absProjectPath, templateData); err != nil {
			log.Fatalf("Error generating project structure: %v", err)
		}

		// 3. Go mod init (already handled by go.mod.tmpl, but we need to ensure it's a valid module)
		// The go.mod.tmpl will create the go.mod file.
		// We'll run `go mod tidy` later to fetch dependencies.

		// 4. Go get Gin (and other initial dependencies if any)
		// Dependencies are listed in go.mod.tmpl, so `go mod tidy` will fetch them.
		fmt.Println("Fetching Gin and other initial dependencies (go mod tidy)...")
		if err := utils.RunCommand(absProjectPath, "go", "mod", "tidy"); err != nil {
			log.Printf("Warning: 'go mod tidy' failed after initial setup: %v. Please run it manually.", err)
		}

		fmt.Println("\nProject", projectName, "created successfully!")
		fmt.Println("To get started:")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Println("  # Review configs/config.yaml")
		fmt.Println("  go run cmd/server/main.go")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringVarP(&modulePath, "module", "m", "", "Go module path (e.g., github.com/user/project)")
}
