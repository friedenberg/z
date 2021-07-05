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
			return z.IndexData.Areas
		},
		ValueSetFunc: func(z *Zettel, v []string) {
			z.IndexData.Areas = v
		},
	}
}

func GetMetadataFieldReadWriterProjects() MetadataFieldReadWriterArray {
	return MetadataFieldReadWriterArray{
		ValueGetFunc: func(z *Zettel) []string {
			return z.IndexData.Projects
		},
		ValueSetFunc: func(z *Zettel, v []string) {
			z.IndexData.Projects = v
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
