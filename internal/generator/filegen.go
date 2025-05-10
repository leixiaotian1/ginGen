package generator

import (
	"fmt"
	"ginGen/internal/utils"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateData holds data for template execution.
type TemplateData struct {
	ProjectName string
	ModulePath  string
	// Add other fields as needed for specific templates
}

// CreateFileFromTemplate generates a file from an embedded template.
func CreateFileFromTemplate(fsys fs.FS, templatePath, outputPath string, data interface{}) error {
	// Ensure the template path uses forward slashes, as expected by embed.FS
	fsTemplatePath := strings.ReplaceAll(templatePath, "\\", "/")

	// Read template content from embedded FS
	templateContent, err := fs.ReadFile(fsys, fsTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", fsTemplatePath, err)
	}

	// Create parent directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	// Parse template
	tmplName := filepath.Base(templatePath)
	tmpl, err := template.New(tmplName).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", tmplName, err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template for %s: %w", outputPath, err)
	}

	fmt.Printf("Generated file: %s\n", outputPath)
	return nil
}

// GenerateProjectStructure creates the initial project files and directories.
func GenerateProjectStructure(projectPath string, data TemplateData) error {
	// Define files to be generated with their template paths and output paths
	filesToGenerate := []struct {
		templatePath string
		outputPath   string
	}{
		{"templates/newproject/main.go.tmpl", "cmd/server/main.go"}, // Changed to cmd/server/main.go
		{"templates/newproject/go.mod.tmpl", "go.mod"},              // Generate go.mod from template
		{"templates/newproject/internal/router/router.go.tmpl", "internal/router/router.go"},
		{"templates/newproject/internal/config/config.go.tmpl", "internal/config/config.go"},
		{"templates/newproject/internal/handler/hello.go.tmpl", "internal/handler/hello.go"},
		{"templates/newproject/configs/config.yaml.tmpl", "configs/config.yaml"},
		{"templates/newproject/dotgitignore.tmpl", ".gitignore"},
		{"templates/newproject/README.md.tmpl", "README.md"},
	}

	for _, f := range filesToGenerate {
		fullOutputPath := filepath.Join(projectPath, f.outputPath)
		err := CreateFileFromTemplate(AllTemplatesFS, f.templatePath, fullOutputPath, data)
		if err != nil {
			return err
		}
	}

	// Create empty directories
	dirsToCreate := []string{
		"internal/service",
		"internal/model",
		"internal/middleware",
		"pkg", // Optional but good to have
	}
	return utils.CreateDirs(projectPath, dirsToCreate...)
}
