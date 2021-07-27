package lib

func GetAlfredFormatDefault() ZettelAlfredItemFormat {
	return ZettelAlfredItemFormat{
		Title:        MakePrintfFormatFunc("%d"),
		Arg:          MakePrintfFormatFunc("%p"),
		Subtitle:     MakePrintfFormatFunc("%w, %t"),
		QuicklookUrl: MakePrintfFormatFunc("%f"),
		IconType:     MakePrintfFormatFunc("file"),
		IconPath:     MakePrintfFormatFunc("%f"),
		// Text:         MakePrintfFormatter("%b"),
		// Match        Formatter
	}
}
