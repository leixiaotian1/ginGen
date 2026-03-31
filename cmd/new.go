package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/leixiaotian1/ginGen/internal/feature"
	"github.com/leixiaotian1/ginGen/internal/generator"
	"github.com/leixiaotian1/ginGen/internal/preset"
	"github.com/leixiaotian1/ginGen/internal/utils"

	"github.com/spf13/cobra"
)

var modulePath string
var newFeatures string
var newPreset string
var newTemplateRoot string

var newCmd = &cobra.Command{
	Use:   "new <project_name>",
	Short: "Create a new Gin project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		targetModulePath := modulePath

		if targetModulePath == "" {
			if d := preset.DefaultModuleOrEmpty(); d != "" {
				targetModulePath = d
				fmt.Printf("Using defaultModule from .gingen.yaml: %s\n", targetModulePath)
			} else {
				targetModulePath = projectName
				fmt.Printf("Module path not specified, defaulting to: %s\n", targetModulePath)
				fmt.Println("You can set defaultModule in .gingen.yaml or use --module, e.g. github.com/youruser/yourproject")
			}
		}

		tmplRoot := newTemplateRoot
		if tmplRoot == "" {
			if p, err := preset.LoadMerged(); err == nil && p.TemplateRoot != "" {
				tmplRoot = p.TemplateRoot
			}
		}

		fmt.Printf("Creating new Gin project: %s (Module: %s)\n", projectName, targetModulePath)

		if err := os.MkdirAll(projectName, 0755); err != nil {
			log.Fatalf("Error creating project directory %s: %v", projectName, err)
		}
		fmt.Printf("Created project directory: %s\n", projectName)

		absProjectPath, err := filepath.Abs(projectName)
		if err != nil {
			log.Fatalf("Error getting absolute path for %s: %v", projectName, err)
		}

		templateData := generator.TemplateData{
			ProjectName: projectName,
			ModulePath:  targetModulePath,
		}
		if err := generator.GenerateProjectWithOptions(absProjectPath, templateData, generator.ProjectGenOptions{
			Quiet:        false,
			TemplateRoot: tmplRoot,
		}); err != nil {
			log.Fatalf("Error generating project structure: %v", err)
		}

		fmt.Println("Fetching Gin and other initial dependencies (go mod tidy)...")
		if err := utils.RunCommand(absProjectPath, "go", "mod", "tidy"); err != nil {
			log.Printf("Warning: 'go mod tidy' failed after initial setup: %v. Please run it manually.", err)
		}

		featureQueue := collectNewProjectFeatures(newFeatures, newPreset)
		if len(featureQueue) > 0 {
			seen := make(map[string]bool)
			for _, raw := range featureQueue {
				raw = strings.TrimSpace(raw)
				if raw == "" {
					continue
				}
				id, err := feature.Normalize(raw)
				if err != nil {
					log.Fatalf("Unknown feature %q: %v", raw, err)
				}
				if seen[id] {
					continue
				}
				seen[id] = true
				ctx := &feature.Context{
					ProjectPath: absProjectPath,
					Data:        templateData,
					Opts: feature.ApplyOptions{
						Quiet:        false,
						Force:        false,
						SkipGoGet:    false,
						TemplateRoot: tmplRoot,
					},
				}
				fmt.Printf("Applying feature: %s\n", id)
				if err := feature.Apply(ctx, id); err != nil {
					log.Fatalf("Error applying feature %s: %v", id, err)
				}
			}
			fmt.Println("Running go mod tidy after features...")
			if err := utils.RunCommand(absProjectPath, "go", "mod", "tidy"); err != nil {
				log.Printf("Warning: go mod tidy failed: %v", err)
			}
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
	newCmd.Flags().StringVar(&newFeatures, "features", "", "comma-separated features to add after creation (e.g. mysql,redis,jwt)")
	newCmd.Flags().StringVar(&newPreset, "preset", "", "name of preset.features from merged .gingen.yaml (combined with --features)")
	newCmd.Flags().StringVar(&newTemplateRoot, "template-root", "", "directory overlaying embedded templates (must contain templates/...)")
}

func collectNewProjectFeatures(commaFeatures, presetName string) []string {
	var out []string
	if strings.TrimSpace(commaFeatures) != "" {
		for _, part := range strings.Split(commaFeatures, ",") {
			out = append(out, strings.TrimSpace(part))
		}
	}
	if strings.TrimSpace(presetName) == "" {
		return out
	}
	cfg, err := preset.LoadMerged()
	if err != nil {
		log.Fatalf("load preset config: %v", err)
	}
	p, ok := cfg.Presets[presetName]
	if !ok {
		log.Fatalf("Unknown preset %q (define presets.%s in .gingen.yaml)", presetName, presetName)
	}
	out = append(out, p.Features...)
	return out
}
