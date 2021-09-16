package lib

import (
	"bufio"
	"io"
	"strings"
)

const (
	MetadataStartSequence = "---\n"
	MetadataEndSequence   = "...\n"
)

func (z *Zettel) ReadFrom(r1 io.Reader, includeBody bool) (err error) {
	r := bufio.NewReader(r1)
	sb := strings.Builder{}
	within := false

	for {
		var s string
		s, err = r.ReadString('\n')

		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			return
		}

		if !within && s == MetadataStartSequence {
			within = true
		} else if within && s == MetadataEndSequence {
			err = z.Metadata.Set(sb.String())

			if err != nil || !includeBody {
				return
			}

			sb.Reset()
			within = false
		} else {
			sb.WriteString(s)
		}
	}

	z.Note.Body = strings.TrimSpace(sb.String())

	return
}
