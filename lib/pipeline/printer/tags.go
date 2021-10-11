package printer

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel/metadata"
	"github.com/friedenberg/z/util"
)

type tagCounts struct {
	zettels, files, urls int
}

type Tag struct {
	metadata.ITag
	tagCounts
}

type Tags struct {
	ShouldExpand bool
	tags         map[string]Tag
}

func (p *Tags) Begin() {
	p.tags = make(map[string]Tag)
}

func (p *Tags) PrintZettel(i int, z *lib.Zettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	for _, t := range z.Note.Metadata.Tags() {
		c, _ := p.tags[t.Tag()]

		c.zettels += 1

		if z.Note.Metadata.HasFile() {
			c.files += 1
		}

		if _, ok := z.Note.Metadata.Url(); ok {
			c.urls += 1
		}

		p.tags[t.Tag()] = c
	}
}

func (p *Tags) End() {
	util.StdPrinterOut(`{"items":[`)

	needsComma := false

	for t, c := range p.tags {
		if needsComma {
			util.StdPrinterOut(",")
		}

		item := alfredItemFromTag(t, c)
		j, err := lib.GenerateAlfredItemsJson([]lib.AlfredItem{item})

		if err != nil {
			//TODO-P2 handle error
			continue
		}

		util.StdPrinterOut(j)
		needsComma = true
	}

	util.StdPrinterOut(`]}`)
}
