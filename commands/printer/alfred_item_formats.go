package printer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

func alfredItemFromZettelBase(z *lib.Zettel) (i lib.AlfredItem) {
	id := strconv.FormatInt(z.Id, 10)
	if len(z.Metadata.Description) > 0 {
		i.Title = z.Metadata.Description
	} else {
		i.Title = z.Format("%w")
	}

	i.Arg = z.Path
	i.Uid = "z." + id
	i.ItemType = "file:skipcheck"

	if len(z.Metadata.Tags) > 0 {
		i.Subtitle = z.Format("%t")
	} else {
		i.Subtitle = z.Format("%w")
	}

	i.QuicklookUrl = z.Format("%f")
	i.Match = MakeAlfredMatches(z)

	i.Icon = lib.AlfredItemIcon{
		Type: "fileicon",
	}

	i.Icon.Path = z.Path

	i.Text = lib.AlfredItemText{
		Copy: id,
	}

	return
}

func alfredItemsFromZettelDefault(z *lib.Zettel) (a []lib.AlfredItem) {
	a = append(a, alfredItemFromZettelBase(z))

	return
}

func AlfredItemsFromZettelFiles(z *lib.Zettel) (a []lib.AlfredItem) {
	i := alfredItemFromZettelBase(z)
	i.Icon.Path = z.FilePath()
	i.Arg = z.FilePath()
	i.Uid = i.Uid + ".file"
	i.Match = i.Match + "i-f"
	a = append(a, i)

	return
}

func AlfredItemsFromZettelUrls(z *lib.Zettel) (a []lib.AlfredItem) {
	i := alfredItemFromZettelBase(z)
	//TODO set to url icon
	// i.Icon.Path = z.FilePath()
	i.Arg = z.Metadata.Url
	i.Title = z.Metadata.Url
	i.Uid = i.Uid + ".url"
	i.Match = i.Match + "i-u"
	a = append(a, i)

	return
}

func AlfredItemsFromZettelAll(z *lib.Zettel) (a []lib.AlfredItem) {
	a = append(a, alfredItemFromZettelBase(z))

	if z.HasFile() {
		a = append(a, AlfredItemsFromZettelFiles(z)...)
	}

	if z.HasUrl() {
		a = append(a, AlfredItemsFromZettelUrls(z)...)
	}

	return
}

func AlfredItemsFromZettelSnippets(z *lib.Zettel) (a []lib.AlfredItem) {
	i := alfredItemFromZettelBase(z)
	//TODO move body normalization to dedicated function
	i.Title = strings.ReplaceAll(z.Body, "\n", " ")
	i.Subtitle = z.Format("%d, %t")
	a = append(a, i)
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
