package lib

func (z *Zettel) HasUrl() bool {
	return z.Metadata.Url != ""
}
