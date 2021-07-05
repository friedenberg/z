package lib

func (z *Zettel) HasFile() bool {
	return z.IndexData.File != ""
}
