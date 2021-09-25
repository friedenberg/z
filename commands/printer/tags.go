package printer

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util"
)

type tagCounts struct {
	zettels, files, urls int
}

type Tags struct {
	ShouldExpand bool
	tags         map[string]tagCounts
}

func (p *Tags) Begin() {
	p.tags = make(map[string]tagCounts)
}

func (p *Tags) PrintZettel(i int, z *lib.KastenZettel, errIn error) {
	if errIn != nil {
		util.StdPrinterErr(errIn)
		return
	}

	tags := z.Metadata.Tags

	if p.ShouldExpand {
		tags = z.Metadata.ExpandedTags
	}

	for _, t := range tags {
		var c tagCounts

		c, _ = p.tags[t]

		c.zettels += 1

		if z.HasFile() {
			c.files += 1
		}

		if z.HasUrl() {
			c.urls += 1
		}

		p.tags[t] = c
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
		//TODO handle error
		j, _ := lib.GenerateAlfredItemsJson([]lib.AlfredItem{item})

		util.StdPrinterOut(j)
		needsComma = true
	}

	util.StdPrinterOut(`]}`)
}
