package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunCommand executes a shell command and prints its output.
func RunCommand(workingDir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if workingDir != "" {
		cmd.Dir = workingDir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("Executing: %s %s (in %s)\n", name, strings.Join(args, " "), workingDir)
	return cmd.Run()
}

// CreateDirs creates necessary directories.
func CreateDirs(basePath string, dirs ...string) error {
	for _, dir := range dirs {
		fullPath := filepath.Join(basePath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", fullPath, err)
		}
		fmt.Printf("Created directory: %s\n", fullPath)
	}
	return nil
}

// GetModulePathFromGoMod reads the module path from a go.mod file.
func GetModulePathFromGoMod(projectPath string) (string, error) {
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("could not read go.mod at %s: %w", goModPath, err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	return "", fmt.Errorf("module directive not found in %s", goModPath)
}
