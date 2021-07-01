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

func GetMetadataFieldReadWriterAreas() MetadataFieldReadWriterArray {
	return MetadataFieldReadWriterArray{
		ValueGetFunc: func(z *Zettel) []string {
			return z.Metadata.Areas
		},
		ValueSetFunc: func(z *Zettel, v []string) {
			z.Metadata.Areas = v
		},
	}
}

func GetMetadataFieldReadWriterProjects() MetadataFieldReadWriterArray {
	return MetadataFieldReadWriterArray{
		ValueGetFunc: func(z *Zettel) []string {
			return z.Metadata.Projects
		},
		ValueSetFunc: func(z *Zettel, v []string) {
			z.Metadata.Projects = v
		},
	}
}

func GetMetadataFieldReadWriterTags() MetadataFieldReadWriterArray {
	return MetadataFieldReadWriterArray{
		ValueGetFunc: func(z *Zettel) []string {
			return z.Metadata.Tags
		},
		ValueSetFunc: func(z *Zettel, v []string) {
			z.Metadata.Tags = v
		},
	}
}
