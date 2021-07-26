package lib

import "path"

func (z *Zettel) HasFile() bool {
	return z.IndexData.File != ""
}

func (z *Zettel) FilePath() string {
	if !z.HasFile() {
		return ""
	}

	return path.Join(z.Env.BasePath, z.IndexData.File)
}
