package writer

import (
	"io"
	"strings"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/util/stdprinter"
)

type Formatter struct {
	formats      []string
	excludeEmpty bool
}

func MakeFormatter(format string) *Formatter {
	formats := strings.Split(format, "%")
	return &Formatter{formats: formats}
}

func (f *Formatter) SetExcludeEmpty() {
	stdprinter.Debug("set exclude empty")
	f.excludeEmpty = true
}

func (f Formatter) WriteZettel(w io.Writer, i int, z *zettel.Zettel) {
	s := f.Format(z)

	if strings.TrimSpace(s) == "" && f.excludeEmpty {
		return
	}

	_, err := io.WriteString(w, s)
	stdprinter.PanicIfError(err)
}

func FormatZettel(z *zettel.Zettel, format string) string {
	f := MakeFormatter(format)
	return f.Format(z)
}

func (f Formatter) Format(z *zettel.Zettel) string {
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
			sb.WriteString(z.Metadata.Description())
		case 'f':
			if f, ok := z.Note.Metadata.LocalFile(); ok {
				sb.WriteString(f.FilePath(z.ZUmwelt.Dir()))
			}
		case 'p':
			sb.WriteString(z.Path)
		case 't':
			sb.WriteString(strings.Join(z.Metadata.TagStrings(), ", "))
		case 'u':
			if u, ok := z.Metadata.Url(); ok {
				sb.WriteString(u.String())
			}
		case 'w':
			t, err := lib.TimeFromPath(z.Path)

			if err != nil {
				panic(err)
			}

			day := t.Format("2006-01-02")
			sb.WriteString(day)
		case 'z':
			sb.WriteString(z.Id.String())
		default:
			didConsume = false
		}

		if didConsume {
			f = f[1:]
		}

		if len(f) > 0 {
			sb.WriteString(f)
		}
	}

	return sb.String()
}
