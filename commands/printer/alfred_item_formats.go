package printer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

//TODO refactor default into base
func alfredItemFromZettelBase(z *lib.Zettel) (i lib.AlfredItem) {
	return
}

func alfredItemFromZettelDefault(z *lib.Zettel) (i lib.AlfredItem) {
	id := strconv.FormatInt(z.Id, 10)
	i.Title = z.IndexData.Description
	i.Arg = z.Path
	i.Uid = "z." + id

	if len(z.IndexData.Tags) > 0 {
		i.Subtitle = z.Format("%t")
	} else {
		i.Subtitle = z.Format("%w")
	}

	i.QuicklookUrl = z.Format("%f")
	i.Match = MakeAlfredMatches(z)

	i.Icon = lib.AlfredItemIcon{
		Type: "fileicon",
	}

	if z.HasFile() {
		i.Icon.Path = z.FilePath()
	} else {
		i.Icon.Path = z.Path
	}

	i.Text = lib.AlfredItemText{
		Copy: id,
	}

	return
}

func alfredItemFromZettelSnippet(z *lib.Zettel) (i lib.AlfredItem) {
	i = alfredItemFromZettelDefault(z)
	i.Title = strings.ReplaceAll(z.Data.Body, "\n", " ")
	i.Subtitle = z.Format("%d, %t")
	return
}

func alfredItemFromTag(t string, counts tagCounts) (i lib.AlfredItem) {
	i.Title = t
	i.Arg = t
	i.Uid = "z." + t

	sb := &strings.Builder{}

	if counts.zettels == 1 {
		sb.WriteString("1 zettel")
	} else {
		sb.WriteString(fmt.Sprintf("%d zettels", counts.zettels))
	}

	addCount := func(name string, c int) {
		if c == 1 {
			sb.WriteString(fmt.Sprintf(", 1 %s", name))
		} else if c > 1 {
			sb.WriteString(fmt.Sprintf(", %d %ss", c, name))
		}
	}

	addCount("file", counts.files)
	addCount("url", counts.urls)

	i.Subtitle = sb.String()

	sb = &strings.Builder{}

	for _, m := range util.ExpandTags(t) {
		sb.WriteString(m)
		sb.WriteString(" ")
	}

	if counts.files > 0 {
		sb.WriteString("h-f")
		sb.WriteString(" ")
	}

	if counts.urls > 0 {
		sb.WriteString("h-u")
		sb.WriteString(" ")
	}

	i.Match = sb.String()

	return
}
