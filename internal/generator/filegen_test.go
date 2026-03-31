package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateProjectStructureQuiet(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	name := "testproj"
	root := filepath.Join(dir, name)
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}
	data := TemplateData{ProjectName: name, ModulePath: "example.com/" + name}
	if err := GenerateProjectStructureQuiet(root, data, true); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{
		"go.mod",
		"cmd/server/main.go",
		"internal/router/router.go",
		"configs/config.yaml",
	} {
		p := filepath.Join(root, rel)
		if _, err := os.Stat(p); err != nil {
			t.Errorf("missing %s: %v", rel, err)
		}
	}
}
