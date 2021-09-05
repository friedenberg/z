package lib

import (
	"path"
)

func (z *Zettel) HasFile() bool {
	return z.Metadata.File != ""
}

func (z *Zettel) FilePath() string {
	if !z.HasFile() {
		return ""
	}

	return path.Join(z.Kasten.BasePath, z.Metadata.File)
}
