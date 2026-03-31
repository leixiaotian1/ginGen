package preset

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const filename = ".gingen.yaml"

// Config is user/project-level ginGen configuration (preset-style).
type Config struct {
	DefaultModule string            `yaml:"defaultModule"`
	TemplateRoot  string            `yaml:"templateRoot"`
	Presets       map[string]Preset `yaml:"presets"`
}

// Preset is a named list of features (for future CLI e.g. `add preset <name>`).
type Preset struct {
	Features []string `yaml:"features"`
}

// LoadMerged loads ~/.gingen.yaml and <cwd>/.gingen.yaml; later files override earlier.
func LoadMerged() (Config, error) {
	var out Config
	home, err := os.UserHomeDir()
	if err == nil {
		p := filepath.Join(home, filename)
		if c, err := loadFile(p); err == nil {
			out = merge(out, c)
		}
	}
	if wd, err := os.Getwd(); err == nil {
		p := filepath.Join(wd, filename)
		if c, err := loadFile(p); err == nil {
			out = merge(out, c)
		}
	}
	return out, nil
}

func loadFile(path string) (Config, error) {
	var c Config
	b, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}
	if err := yaml.Unmarshal(b, &c); err != nil {
		return c, err
	}
	return c, nil
}

func merge(a, b Config) Config {
	if b.DefaultModule != "" {
		a.DefaultModule = b.DefaultModule
	}
	if b.TemplateRoot != "" {
		a.TemplateRoot = b.TemplateRoot
	}
	if len(b.Presets) > 0 {
		if a.Presets == nil {
			a.Presets = map[string]Preset{}
		}
		for k, v := range b.Presets {
			a.Presets[k] = v
		}
	}
	return a
}

// DefaultModuleOrEmpty returns merged default module path, or empty if unset.
func DefaultModuleOrEmpty() string {
	c, err := LoadMerged()
	if err != nil {
		return ""
	}
	return c.DefaultModule
}
