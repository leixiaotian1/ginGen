package feature

import (
	"io/fs"
	"os"

	"github.com/leixiaotian1/ginGen/internal/generator"
)

// ApplyOptions controls feature application (CLI vs Web, strictness, tests).
type ApplyOptions struct {
	Quiet        bool
	Force        bool   // if true, continue after go get / tidy failures (warn only where applicable)
	SkipGoGet    bool   // for unit tests without network
	TemplateRoot string // if set, overlay files from <root>/templates/... over embed
}

// Context carries one apply operation.
type Context struct {
	ProjectPath string
	Data        generator.TemplateData
	Opts        ApplyOptions
}

// TemplateFS returns embed FS or overlay with user template root.
func (c *Context) TemplateFS() fs.FS {
	base := generator.AllTemplatesFS
	if c.Opts.TemplateRoot == "" {
		return base
	}
	return generator.OverlayFS{
		Primary:   os.DirFS(c.Opts.TemplateRoot),
		Secondary: base,
	}
}
