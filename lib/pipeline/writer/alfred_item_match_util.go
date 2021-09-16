package writer

import (
	"strings"
	"time"

	"github.com/friedenberg/z/lib"
	"golang.org/x/xerrors"
)

func MakeAlfredMatches(z *lib.Zettel) string {
	m := z.Metadata
	sb := &strings.Builder{}

	addMatch := func(s string) {
		sb.WriteString(s)
		sb.WriteString(" ")
	}

	t, err := lib.TimeFromPath(z.Path)

	if err != nil {
		panic(xerrors.Errorf("make alfred match field: %w", err))
	}

	addMatch(m.Description())

	for _, t1 := range m.SearchMatchTagStrings() {
		addMatch(t1)
	}

	day := t.Format("2006-01-02")

	addMatch("w-" + day)

	today := time.Now()

	if today.Format("2006-01-02") == day {
		addMatch("w-today")
	}

	return sb.String()
}
