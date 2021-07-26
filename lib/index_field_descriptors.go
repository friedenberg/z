package lib

type MetadataFieldReadWriterArray struct {
	ValueGetFunc func(z *Zettel) []string
	ValueSetFunc func(z *Zettel, v []string)
}

func GetMetadataFieldReadWriterNull() MetadataFieldReadWriterArray {
	return MetadataFieldReadWriterArray{
		ValueGetFunc: func(_ *Zettel) []string {
			return []string{}
		},
		ValueSetFunc: func(_ *Zettel, _ []string) {
		},
	}
}

func GetMetadataFieldReadWriterTags() MetadataFieldReadWriterArray {
	return MetadataFieldReadWriterArray{
		ValueGetFunc: func(z *Zettel) []string {
			return z.IndexData.Tags
		},
		ValueSetFunc: func(z *Zettel, v []string) {
			z.IndexData.Tags = v
		},
	}
}
