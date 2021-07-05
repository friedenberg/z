package lib

func (z *Zettel) HasUrl() bool {
	return z.IndexData.Url != ""
}
