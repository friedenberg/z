package writer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
)

func AlfredItemFromZettelBase(z *zettel.Zettel) (i lib.AlfredItem) {
	id := strconv.FormatInt(z.Id, 10)
	if len(z.Note.Metadata.Description()) > 0 {
		i.Title = z.Note.Metadata.Description()
	} else {
		i.Title = FormatZettel(z, "%w")
	}

	i.Arg = z.Path
	i.Uid = "z." + id
	i.ItemType = "file:skipcheck"

	if z.Note.Metadata.StringTags().Len() > 0 {
		i.Subtitle = FormatZettel(z, "%t")
	} else {
		i.Subtitle = FormatZettel(z, "%w")
	}

	i.QuicklookUrl = FormatZettel(z, "%f")
	mb := MakeMatchBuilder()
	i.Match = mb.Zettel(z)

	i.Icon = lib.AlfredItemIcon{
		Type: "fileicon",
	}

	i.Icon.Path = z.Path

	i.Text = lib.AlfredItemText{
		Copy: id,
	}

	return
}

func alfredItemsFromZettelDefault(z *zettel.Zettel) (a []lib.AlfredItem) {
	a = append(a, AlfredItemFromZettelBase(z))

	return
}

func AlfredItemsFromZettelFiles(z *zettel.Zettel) (a []lib.AlfredItem) {
	f, ok := z.Note.Metadata.LocalFile()

	if !ok {
		return
	}

	i := AlfredItemFromZettelBase(z)
	i.Icon.Path = f.FilePath(z.ZUmwelt.Dir())
	i.Uid = i.Uid + ".file"
	i.Match = i.Match + "i-f"
	i.Arg = fmt.Sprintf("-actions open-files %s", z.Path)
	a = append(a, i)

	return
}

func AlfredItemsFromZettelUrls(z *zettel.Zettel) (a []lib.AlfredItem) {
	u, ok := z.Note.Metadata.Url()

	if !ok {
		return
	}

	i := AlfredItemFromZettelBase(z)
	//TODO-P2 set to url icon
	// i.Icon.Path = z.FilePath()
	i.Arg = fmt.Sprintf("-actions open-urls %s", z.Path)

	//TODO-P4 move to tags
	i.Title = fmt.Sprintf("%s: %s", u.Host, z.Note.Metadata.Description())
	i.Uid = i.Uid + ".url"
	i.Match = i.Match + "i-u"
	a = append(a, i)

	return
}

func AlfredItemsFromZettelAll(z *zettel.Zettel) (a []lib.AlfredItem) {
	a = append(a, AlfredItemFromZettelBase(z))

	if z.Note.Metadata.HasFile() {
		a = append(a, AlfredItemsFromZettelFiles(z)...)
	}

	if _, ok := z.Note.Metadata.Url(); ok {
		a = append(a, AlfredItemsFromZettelUrls(z)...)
	}

	return
}

func AlfredItemsFromZettelSnippets(z *zettel.Zettel) (a []lib.AlfredItem) {
	i := AlfredItemFromZettelBase(z)
	//TODO-P3 move body normalization to dedicated function
	i.Title = strings.ReplaceAll(z.Body, "\n", " ")
	i.Subtitle = FormatZettel(z, "%d, %t")
	a = append(a, i)
	return
}

func alfredItemFromTag(t string, t1 Tag) (i lib.AlfredItem) {
	i.Title = t
	i.Arg = t
	i.Uid = "z." + t

	sb := &strings.Builder{}

	if t1.zettels == 1 {
		sb.WriteString("1 zettel")
	} else {
		sb.WriteString(fmt.Sprintf("%d zettels", t1.zettels))
	}

	addCount := func(name string, c int) {
		if c == 1 {
			sb.WriteString(fmt.Sprintf(", 1 %s", name))
		} else if c > 1 {
			sb.WriteString(fmt.Sprintf(", %d %ss", c, name))
		}
	}

	addCount("file", t1.files)
	addCount("url", t1.urls)

	i.Subtitle = sb.String()

	mb := MakeMatchBuilder()

	mb.AddMatch(t1.Tag())
	for _, m := range t1.SearchMatchTags().Strings() {
		mb.AddMatch(m)
	}

	i.Match = mb.String()

	return
}
