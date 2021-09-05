package printer

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

func MakeAlfredMatches(z *lib.Zettel) string {
	//TODO add more variations and match against item format
	//e.g., Project: 2020-zettel -> p:2020-zettel, p:2020, 2020, zettel
	m := z.Metadata
	sb := &strings.Builder{}

	addMatch := func(s string) {
		sb.WriteString(s)
		sb.WriteString(" ")
	}

	t, err := lib.TimeFromPath(z.Path)

	if err != nil {
		panic(fmt.Errorf("make alfred match field: %w", err))
	}

	addMatch(m.Description)

	day := t.Format("2006-01-02")

	addMatch("w-" + day)

	if z.HasUrl() {
		url, err := url.Parse(m.Url)

		if err == nil {
			addMatch("d-" + url.Hostname())
		}

		addMatch("h-u")
	}

	if z.HasFile() {
		addMatch("h-f")
	}

	today := time.Now()

	if today.Format("2006-01-02") == day {
		addMatch("w-today")
	}

	for _, t := range m.Tags {
		for _, m := range util.ExpandTags(t) {
			addMatch(m)
		}
	}

	return sb.String()
}
