package generator

import (
	"io/fs"
)

// OverlayFS reads from primary first; if Open fails, falls back to secondary.
// Used to layer user templates over embedded defaults.
type OverlayFS struct {
	Primary   fs.FS
	Secondary fs.FS
}

// Open implements fs.FS.
func (o OverlayFS) Open(name string) (fs.File, error) {
	if o.Primary != nil {
		f, err := o.Primary.Open(name)
		if err == nil {
			return f, nil
		}
	}
	if o.Secondary == nil {
		return nil, fs.ErrNotExist
	}
	return o.Secondary.Open(name)
}
