package writer

import (
	"strings"
	"time"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/util/stdprinter"
	"golang.org/x/xerrors"
)

type MatchBuilder struct {
	*strings.Builder
}

func MakeMatchBuilder() MatchBuilder {
	return MatchBuilder{
		Builder: &strings.Builder{},
	}
}

func (mb MatchBuilder) AddMatch(s string) {
	s1 := strings.Split(s, "_")

	for _, s2 := range s1 {
		mb.WriteString(s2)
		mb.WriteString(" ")
	}

	mb.WriteString(s)
	mb.WriteString(" ")
}

func (mb *MatchBuilder) Zettel(z *zettel.Zettel) string {
	m := z.Metadata

	t, err := lib.TimeFromPath(z.Path)

	if err != nil {
		err = xerrors.Errorf("make alfred match field: %w", err)
		stdprinter.PanicIfError(err)
	}

	mb.AddMatch(m.Description())

	for _, t1 := range m.SearchMatchTagStrings() {
		mb.AddMatch(t1)
	}

	day := t.Format("2006-01-02")

	mb.AddMatch("w-" + day)

	today := time.Now()

	if today.Format("2006-01-02") == day {
		mb.AddMatch("w-today")
	}

	return mb.String()
}
