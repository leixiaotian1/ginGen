package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/leixiaotian1/ginGen/internal/feature"
	"github.com/leixiaotian1/ginGen/internal/generator"
	"github.com/leixiaotian1/ginGen/internal/preset"
	"github.com/leixiaotian1/ginGen/internal/utils"

	"github.com/spf13/cobra"
)

var addForce bool
var addTemplateRoot string

var addCmd = &cobra.Command{
	Use:   "add <feature> [project_path]",
	Short: "Add a feature (e.g., mysql) to an existing project",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		featureArg := args[0]
		projectPath := "."
		if len(args) > 1 {
			projectPath = args[1]
		}

		absProjectPath, err := filepath.Abs(projectPath)
		if err != nil {
			log.Fatalf("Error getting absolute path for %s: %v", projectPath, err)
		}

		goModPath := filepath.Join(absProjectPath, "go.mod")
		if _, err := os.Stat(goModPath); os.IsNotExist(err) {
			log.Fatalf("Error: %s is not a Go module (go.mod not found).\n", absProjectPath)
		}

		currentModulePath, err := utils.GetModulePathFromGoMod(absProjectPath)
		if err != nil {
			log.Fatalf("Error reading module path from go.mod: %v", err)
		}

		id, err := feature.Normalize(featureArg)
		if err != nil {
			log.Fatalf("Error: %v. Supported features: %s.\n", err, feature.SupportedListForError())
		}

		tmplRoot := addTemplateRoot
		if tmplRoot == "" {
			if p, err := preset.LoadMerged(); err == nil && p.TemplateRoot != "" {
				tmplRoot = p.TemplateRoot
			}
		}

		ctx := &feature.Context{
			ProjectPath: absProjectPath,
			Data:        generator.TemplateData{ModulePath: currentModulePath},
			Opts: feature.ApplyOptions{
				Quiet:        false,
				Force:        addForce,
				SkipGoGet:    false,
				TemplateRoot: tmplRoot,
			},
		}

		fmt.Printf("Adding feature '%s' to project in '%s'...\n", id, absProjectPath)
		if err := feature.Apply(ctx, id); err != nil {
			if addForce {
				log.Printf("Warning: feature apply error (continuing due to --force): %v\n", err)
			} else {
				log.Fatalf("Error: %v\n", err)
			}
		}

		printPostHints(id)

		fmt.Println("Running 'go mod tidy' to update dependencies...")
		if err := utils.RunCommand(absProjectPath, "go", "mod", "tidy"); err != nil {
			if addForce {
				log.Printf("Warning: 'go mod tidy' failed: %v. Please run it manually.\n", err)
			} else {
				log.Fatalf("'go mod tidy' failed: %v (use --force to continue anyway)\n", err)
			}
		}

		fmt.Printf("\nFeature '%s' added to project %s.\n", id, absProjectPath)
		fmt.Println("Review generated files; wire routes/clients in main if needed.")
	},
}

func printPostHints(id string) {
	switch id {
	case "mysql", "redis", "kafka":
		fmt.Println("\n--- Next steps ---")
		fmt.Println("Configs were merged into configs/config.yaml and internal/config/config.go where possible.")
		fmt.Println("Initialize clients in cmd/server/main.go and register routes as needed.")
		fmt.Println("------------------")
	case "swagger":
		fmt.Println("\nNote: Run 'swag init -g cmd/server/main.go' to generate Swagger docs.")
	}
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolVar(&addForce, "force", false, "continue even if go get or go mod tidy fails")
	addCmd.Flags().StringVar(&addTemplateRoot, "template-root", "", "directory overlaying embedded templates (must contain templates/...)")
}
