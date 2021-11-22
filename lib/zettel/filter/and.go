package filter

import (
	"fmt"
	"strings"

	"github.com/friedenberg/z/lib/zettel"
)

type And Filters

func MakeAnd(fs ...Filter) (f *And) {
	f1 := And(fs)
	f = &f1
	return
}

func (f And) FilterZettel(i int, z *zettel.Zettel) bool {
	for _, f1 := range f {
		if f1 == nil {
			continue
		}

		if !f1.FilterZettel(i, z) {
			return false
		}
	}

	return true
}

func (f And) String() string {
	return fmt.Sprintf("%#v", f)
}

func (f *And) Set(v string) (err error) {
	if v != "" {
		var a Filter

		if strings.Contains(v, " ") {
			a = Description(v)
		} else {
			a = Tag(v)
		}

		*f = append(*f, a)
	}

	return
}
