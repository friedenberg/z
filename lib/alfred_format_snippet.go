package lib

import "strings"

func GetAlfredFormatSnippet() ZettelAlfredItemFormat {
	return ZettelAlfredItemFormat{
		Title: func(z *Zettel) (s string) {
			s = MakePrintfFormatFunc("%b")(z)
			s = strings.ReplaceAll(s, "\n", " ")
			return
		},
		Arg:          MakePrintfFormatFunc("%p"),
		Subtitle:     MakePrintfFormatFunc("%d, %t"),
		QuicklookUrl: MakePrintfFormatFunc("%f"),
		IconType:     MakePrintfFormatFunc("file"),
		IconPath:     MakePrintfFormatFunc("%f"),
		// Text:         MakePrintfFormatter("%b"),
		// Match        Formatter
	}
}
