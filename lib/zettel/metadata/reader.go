package metadata

import (
	"bufio"
	"io"
	"strings"
)

const (
	MetadataStartSequence = "---\n"
	MetadataEndSequence   = "...\n"
)

func ReadYAMLHeader(r1 io.Reader) (h string, err error) {
	r := bufio.NewReader(r1)
	sb := strings.Builder{}
	within := false

	for {
		var some_string string
		some_string, err = r.ReadString('\n')

		if err == io.EOF {
			err = nil
			return
		}

		if err != nil {
			return
		}

		if !within && some_string == MetadataStartSequence {
			within = true
		} else if within && some_string == MetadataEndSequence {
			break
		} else if within {
			sb.WriteString(some_string)
		}
	}

	h = sb.String()

	return
}
