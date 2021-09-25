package lib

func (z *KastenZettel) HasUrl() bool {
	return z.Metadata.Url != ""
}
