package lib

import (
	"path"
)

func (z *KastenZettel) HasFile() bool {
	return z.Metadata.File != ""
}

func (z *KastenZettel) FilePath() string {
	if !z.HasFile() {
		return ""
	}

	return path.Join(z.Kasten.Local.BasePath, z.Metadata.File)
}
