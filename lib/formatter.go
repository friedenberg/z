package lib

import (
	"strconv"
	"strings"
)

type FormatFunc func(*KastenZettel) string

type Formatter interface {
	Format(*KastenZettel) string
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

func (f printfFormatter) Format(z *KastenZettel) string {
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
		case 'b':
			sb.WriteString(strings.TrimSpace(z.Body))
			// sb.WriteString(strings.ReplaceAll(z.Data.Body, "%", "%%"))
		case 'd':
			sb.WriteString(z.Metadata.Description)
		case 'f':
			sb.WriteString(z.FilePath())
		case 'p':
			sb.WriteString(z.Path)
		case 't':
			sb.WriteString(strings.Join(z.Metadata.Tags, ", "))
		case 'u':
			sb.WriteString(z.Metadata.Url)
		case 'w':
			sb.WriteString(z.Metadata.Date)
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

func (z *KastenZettel) Format(f string) string {
	return MakePrintfFormatter(f).Format(z)
}
