package lib

type MetadataFieldReadWriterArray struct {
	ValueGetFunc func(z *KastenZettel) []string
	ValueSetFunc func(z *KastenZettel, v []string)
}

func GetMetadataFieldReadWriterNull() MetadataFieldReadWriterArray {
	return MetadataFieldReadWriterArray{
		ValueGetFunc: func(_ *KastenZettel) []string {
			return []string{}
		},
		ValueSetFunc: func(_ *KastenZettel, _ []string) {
		},
	}
}

func GetMetadataFieldReadWriterTags() MetadataFieldReadWriterArray {
	return MetadataFieldReadWriterArray{
		ValueGetFunc: func(z *KastenZettel) []string {
			return z.Metadata.Tags
		},
		ValueSetFunc: func(z *KastenZettel, v []string) {
			z.Metadata.Tags = v
		},
	}
}
