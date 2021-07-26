package lib

import (
	"strconv"
	"strings"
)

type FormatFunc func(*Zettel) string

type Formatter interface {
	Format(*Zettel) string
}

type printfFormatter struct {
	formats []string
}

func MakePrintfFormatter(format string) Formatter {
	formats := strings.Split(format, "%")
	return &printfFormatter{formats: formats}
}

func MakePrintfFormatFunc(format string) FormatFunc {
	formats := strings.Split(format, "%")
	f := printfFormatter{formats: formats}
	return f.Format
}

func (f printfFormatter) Format(z *Zettel) string {
	sb := &strings.Builder{}

	lastLoopWasEmpty := false

	if len(f.formats) == 1 {
		return f.formats[0]
	}

	for _, f := range f.formats {
		if len(f) == 0 && lastLoopWasEmpty {
			sb.WriteString("%")
			lastLoopWasEmpty = false
			continue
		} else if len(f) == 0 {
			lastLoopWasEmpty = true
			continue
		}
		didConsume := true

		switch f[0] {
		case 'a':
			sb.WriteString(strings.Join(z.IndexData.Areas, ","))
		case 'b':
			sb.WriteString(z.Data.Body)
			// sb.WriteString(strings.ReplaceAll(z.Data.Body, "%", "%%"))
		case 'd':
			sb.WriteString(z.IndexData.Description)
		case 'f':
			sb.WriteString(z.FilePath())
		case 'p':
			sb.WriteString(z.Path)
		case 'r':
			sb.WriteString(strings.Join(z.IndexData.Projects, ","))
		case 't':
			sb.WriteString(strings.Join(z.IndexData.Tags, ","))
		case 'u':
			sb.WriteString(z.IndexData.Url)
		case 'w':
			sb.WriteString(z.IndexData.Date)
		case 'z':
			sb.WriteString(strconv.FormatInt(z.Id, 10))
		default:
			didConsume = false
		}

		if didConsume {
			f = f[1:]
		}

		if len(f) > 0 {
			//TODO backslash escapes / printf style escapes
			sb.WriteString(f)
		}
	}

	return sb.String()
}