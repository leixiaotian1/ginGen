package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/leixiaotian1/ginGen/internal/generator"
	"github.com/leixiaotian1/ginGen/internal/utils"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <feature> [project_path]",
	Short: "Add a feature (e.g., mysql) to an existing project",
	Args:  cobra.MinimumNArgs(1), // feature is required
	Run: func(cmd *cobra.Command, args []string) {
		feature := args[0]
		projectPath := "." // Default to current directory
		if len(args) > 1 {
			projectPath = args[1]
		}

		absProjectPath, err := filepath.Abs(projectPath)
		if err != nil {
			log.Fatalf("Error getting absolute path for %s: %v", projectPath, err)
		}

		// 1. Verify projectPath is a Go module
		goModPath := filepath.Join(absProjectPath, "go.mod")
		if _, err := os.Stat(goModPath); os.IsNotExist(err) {
			log.Fatalf("Error: %s is not a Go module (go.mod not found).\n", absProjectPath)
		}

		// Get module path from existing go.mod for template data
		currentModulePath, err := utils.GetModulePathFromGoMod(absProjectPath)
		if err != nil {
			log.Fatalf("Error reading module path from go.mod: %v", err)
		}
		templateData := generator.TemplateData{ModulePath: currentModulePath}

		fmt.Printf("Adding feature '%s' to project in '%s'...\n", feature, absProjectPath)

		switch feature {
		case "mysql", "gorm": // Treat them as the same for MVP
			addMySQLGormFeature(absProjectPath, templateData)
		// case "redis":
		//  addRedisFeature(absProjectPath, templateData) // For future
		default:
			log.Fatalf("Error: Unknown feature '%s'. Supported features for MVP: mysql (or gorm).\n", feature)
		}

		fmt.Println("Running 'go mod tidy' to update dependencies...")
		if err := utils.RunCommand(absProjectPath, "go", "mod", "tidy"); err != nil {
			log.Printf("Warning: 'go mod tidy' failed: %v. Please run it manually.\n", err)
		}

		fmt.Printf("\nFeature '%s' added to project %s.\n", feature, absProjectPath)
		fmt.Println("Please review the generated files and update your main configuration and application logic as needed.")
	},
}

func addMySQLGormFeature(projectPath string, data generator.TemplateData) {
	fmt.Println("Adding MySQL (with GORM) support...")

	// 1. Go get GORM and MySQL driver
	pkgs := []string{
		"gorm.io/gorm",
		"gorm.io/driver/mysql",
		"github.com/go-sql-driver/mysql", // GORM's MySQL driver depends on this
	}
	for _, pkg := range pkgs {
		fmt.Printf("Getting package: %s\n", pkg)
		if err := utils.RunCommand(projectPath, "go", "get", pkg); err != nil {
			log.Printf("Warning: 'go get %s' failed: %v. Please ensure it's added or try manually.\n", pkg, err)
			// For MVP, we'll continue, but in a real app, you might want to stop.
		}
	}

	// 2. Generate GORM client and DB config files from templates
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/addfeature/mysql/config.go.tmpl", "internal/config/db_config.go"},
		{"templates/addfeature/mysql/client.go.tmpl", "internal/clients/gorm.go"},
	}

	// Ensure internal/clients directory exists
	if err := utils.CreateDirs(projectPath, "internal/clients"); err != nil {
		log.Fatalf("Failed to create internal/clients directory: %v", err)
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		err := generator.CreateFileFromTemplate(generator.AllTemplatesFS, f.templatePath, fullOutputPath, data)
		if err != nil {
			log.Fatalf("Error generating file %s: %v", f.outputPath, err)
		}
	}

	// 3. Print message about adding to config.yaml and initializing in main.
	fmt.Println("\n--- Action Required ---")
	fmt.Println("1. Update 'configs/config.yaml' with your MySQL connection details.")
	fmt.Println("   Add a section like this (example):")

	// Print the content of config_entry.yaml.tmpl
	yamlSnippetPath := "templates/addfeature/mysql/config_entry.yaml.tmpl"
	yamlSnippet, err := generator.AllTemplatesFS.ReadFile(yamlSnippetPath)
	if err != nil {
		log.Printf("Warning: could not read YAML snippet template %s: %v\n", yamlSnippetPath, err)
		fmt.Print(`  db:
		mysql:
		  dsn: "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		  max_idle_conns: 10
		  max_open_conns: 100
		  conn_max_lifetime: 3600 # in seconds
	`)
		fmt.Println() // 单独添加一个换行
	} else {
		fmt.Println(string(yamlSnippet))
	}

	fmt.Println("2. Update 'internal/config/config.go':")
	fmt.Println("   - Add 'DBConfig `mapstructure:\"db\"`' to the `Config` struct.")
	fmt.Println("   - Ensure `LoadConfig()` loads this new section.")
	fmt.Println("3. Initialize the GORM client in 'cmd/server/main.go' or your application setup:")
	fmt.Println("   - Example: `db, err := clients.NewGORMClient(appConfig.DB.MySQL)`")
	fmt.Println("   - Pass the `*gorm.DB` instance to your services/handlers as needed.")
	fmt.Println("----------------------")
}

func init() {
	rootCmd.AddCommand(addCmd)
}
