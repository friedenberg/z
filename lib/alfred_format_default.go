package lib

func GetAlfredFormatDefault() ZettelAlfredItemFormat {
	return ZettelAlfredItemFormat{
		Title: MakePrintfFormatFunc("%d"),
		Arg:   MakePrintfFormatFunc("%p"),
		Subtitle: func(z *Zettel) (s string) {
			var f FormatFunc
			if len(z.IndexData.Tags) > 0 {
				f = MakePrintfFormatFunc("%t")
			} else {
				f = MakePrintfFormatFunc("%w")
			}
			s = f(z)
			return
		},
		QuicklookUrl: MakePrintfFormatFunc("%f"),
		IconType:     MakePrintfFormatFunc("file"),
		IconPath:     MakePrintfFormatFunc("%f"),
		// Text:         MakePrintfFormatter("%b"),
		// Match        Formatter
	}
}
