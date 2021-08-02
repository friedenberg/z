package printer

import (
	"strconv"
	"strings"

	"github.com/friedenberg/z/lib"
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
