package feature

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leixiaotian1/ginGen/internal/generator"
)

func TestNormalizeAliases(t *testing.T) {
	t.Parallel()
	cases := map[string]string{
		"mysql":      "mysql",
		"gorm":       "mysql",
		"GORM":       "mysql",
		"postgres":   "postgres",
		"logging":    "logger",
		"hot-reload": "hotreload",
	}
	for in, want := range cases {
		got, err := Normalize(in)
		if err != nil || got != want {
			t.Errorf("Normalize(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
}

func TestApplyMySQLSkipGoGet(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	name := "demo"
	root := filepath.Join(dir, name)
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}
	data := generator.TemplateData{ProjectName: name, ModulePath: "example.com/demo"}
	if err := generator.GenerateProjectStructureQuiet(root, data, true); err != nil {
		t.Fatal(err)
	}
	ctx := &Context{
		ProjectPath: root,
		Data:        data,
		Opts: ApplyOptions{
			Quiet:     true,
			SkipGoGet: true,
		},
	}
	if err := Apply(ctx, "mysql"); err != nil {
		t.Fatal(err)
	}
	dbCfg := filepath.Join(root, "internal", "config", "db_config.go")
	if _, err := os.Stat(dbCfg); err != nil {
		t.Fatal(err)
	}
}
