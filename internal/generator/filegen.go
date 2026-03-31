package generator

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/leixiaotian1/ginGen/internal/utils"
)

// TemplateData holds data for template execution.
type TemplateData struct {
	ProjectName string
	ModulePath  string
	// Add other fields as needed for specific templates
}

// CreateFileFromTemplate generates a file from an embedded template (prints progress to stdout).
func CreateFileFromTemplate(fsys fs.FS, templatePath, outputPath string, data interface{}) error {
	return CreateFileFromTemplateQuiet(fsys, templatePath, outputPath, data, false)
}

// CreateFileFromTemplateQuiet generates a file from a template; if quiet is true, no stdout progress.
func CreateFileFromTemplateQuiet(fsys fs.FS, templatePath, outputPath string, data interface{}, quiet bool) error {
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

	if !quiet {
		fmt.Printf("Generated file: %s\n", outputPath)
	}
	return nil
}

// GenerateProjectStructure creates the initial project files and directories.
func GenerateProjectStructure(projectPath string, data TemplateData) error {
	return GenerateProjectWithOptions(projectPath, data, ProjectGenOptions{Quiet: false})
}

// ProjectGenOptions configures new-project generation.
type ProjectGenOptions struct {
	Quiet        bool
	TemplateRoot string // optional OS path; overlays embedded templates/addfeature & newproject paths
}

// GenerateProjectWithOptions creates the initial project (optional template overlay).
func GenerateProjectWithOptions(projectPath string, data TemplateData, o ProjectGenOptions) error {
	return generateProjectInternal(projectPath, data, o)
}

// GenerateProjectStructureQuiet creates the initial project; quiet suppresses per-file and mkdir logs.
func GenerateProjectStructureQuiet(projectPath string, data TemplateData, quiet bool) error {
	return generateProjectInternal(projectPath, data, ProjectGenOptions{Quiet: quiet})
}

func templateFSForRoot(templateRoot string) fs.FS {
	if templateRoot == "" {
		return AllTemplatesFS
	}
	return OverlayFS{Primary: os.DirFS(templateRoot), Secondary: AllTemplatesFS}
}

func generateProjectInternal(projectPath string, data TemplateData, o ProjectGenOptions) error {
	fsys := templateFSForRoot(o.TemplateRoot)
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
		err := CreateFileFromTemplateQuiet(fsys, f.templatePath, fullOutputPath, data, o.Quiet)
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
	return utils.CreateDirsQuiet(projectPath, o.Quiet, dirsToCreate...)
}
