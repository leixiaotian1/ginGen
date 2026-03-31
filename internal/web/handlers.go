package web

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leixiaotian1/ginGen/internal/feature"
	"github.com/leixiaotian1/ginGen/internal/generator"
	"github.com/leixiaotian1/ginGen/internal/preset"
	"github.com/leixiaotian1/ginGen/internal/utils"
)

// ProjectRequest represents the request body for project generation
type ProjectRequest struct {
	ProjectName string   `json:"projectName" binding:"required"`
	ModulePath  string   `json:"modulePath" binding:"required"`
	Features    []string `json:"features"`
}

// GetAvailableFeatures returns the list of available features/modules (from registry).
func GetAvailableFeatures(c *gin.Context) {
	out := make([]map[string]interface{}, 0, len(feature.AllSpecs()))
	for _, s := range feature.AllSpecs() {
		if s == nil {
			continue
		}
		out = append(out, map[string]interface{}{
			"id":          s.ID,
			"name":        s.Name,
			"description": s.Description,
			"category":    s.Category,
		})
	}
	c.JSON(200, gin.H{"features": out})
}

// GenerateProject handles project generation request
func GenerateProject(c *gin.Context) {
	var req ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if req.ProjectName == "" {
		c.JSON(400, gin.H{"error": "Project name is required"})
		return
	}
	if req.ModulePath == "" {
		req.ModulePath = req.ProjectName
	}

	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("ginGen_%s_%d", req.ProjectName, time.Now().Unix()))
	defer os.RemoveAll(tempDir)

	projectPath := filepath.Join(tempDir, req.ProjectName)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		log.Printf("Error creating project directory: %v", err)
		c.JSON(500, gin.H{"error": "Failed to create project directory"})
		return
	}

	templateData := generator.TemplateData{
		ProjectName: req.ProjectName,
		ModulePath:  req.ModulePath,
	}

	tmplRoot := ""
	if p, err := preset.LoadMerged(); err == nil {
		tmplRoot = p.TemplateRoot
	}

	if err := generator.GenerateProjectWithOptions(projectPath, templateData, generator.ProjectGenOptions{
		Quiet:        true,
		TemplateRoot: tmplRoot,
	}); err != nil {
		log.Printf("Error generating project structure: %v", err)
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to generate project: %v", err)})
		return
	}

	seen := make(map[string]bool)
	for _, raw := range req.Features {
		id, err := feature.Normalize(raw)
		if err != nil {
			c.JSON(400, gin.H{"error": fmt.Sprintf("Unknown feature %q", raw)})
			return
		}
		if seen[id] {
			continue
		}
		seen[id] = true
		ctx := &feature.Context{
			ProjectPath: projectPath,
			Data:        templateData,
			Opts: feature.ApplyOptions{
				Quiet:        true,
				Force:        false,
				SkipGoGet:    false,
				TemplateRoot: tmplRoot,
			},
		}
		if err := feature.Apply(ctx, id); err != nil {
			log.Printf("Error adding feature %s: %v", id, err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add feature %s: %v", id, err)})
			return
		}
	}

	if err := utils.RunCommandSilent(projectPath, "go", "mod", "tidy"); err != nil {
		log.Printf("Warning: go mod tidy failed: %v", err)
	}

	zipPath := filepath.Join(tempDir, fmt.Sprintf("%s.zip", req.ProjectName))
	if err := createZip(projectPath, zipPath); err != nil {
		log.Printf("Error creating ZIP file: %v", err)
		c.JSON(500, gin.H{"error": "Failed to create project archive"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", req.ProjectName))
	c.Header("Content-Type", "application/zip")
	c.File(zipPath)
}

func createZip(sourceDir, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		relPath = strings.ReplaceAll(relPath, "\\", "/")
		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(zipEntry, file)
		return err
	})
}
