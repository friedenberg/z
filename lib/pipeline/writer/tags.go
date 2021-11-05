package writer

import (
	"io"
	"sync"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/metadata"
	"github.com/friedenberg/z/util/stdprinter"
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
	tagChan      chan metadata.ITag
	wg           *sync.WaitGroup
}

func (p *Tags) Begin(_ io.Writer) {
	p.tags = make(map[string]Tag)
	p.tagChan = make(chan metadata.ITag)
	p.wg = &sync.WaitGroup{}
	p.wg.Add(1)

	go func() {
		defer p.wg.Done()

		for t := range p.tagChan {
			c, _ := p.tags[t.Tag()]
			c.ITag = t

			c.zettels += 1

			//TODO-P3 add back support for file and url counts
			// if z.Note.Metadata.HasFile() {
			// 	c.files += 1
			// }

			// if _, ok := z.Note.Metadata.Url(); ok {
			// 	c.urls += 1
			// }

			p.tags[t.Tag()] = c
		}
	}()
}

func (p *Tags) WriteZettel(_ io.Writer, i int, z *zettel.Zettel) {
	//TODO-P4 check performance of this
	tagsToPush := z.Note.Metadata.StringTags()

	if p.ShouldExpand {
		tagsToPush = z.Note.Metadata.SearchMatchTags()
	}

	for _, t := range tagsToPush.Tags() {
		p.tagChan <- t
	}
}

func (p *Tags) End(w io.Writer) {
	close(p.tagChan)
	p.wg.Wait()

	var err error
	_, err = io.WriteString(w, `{"items":[`)

	stdprinter.PanicIfError(err)

	needsComma := false

	for t, c := range p.tags {
		if needsComma {
			_, err = io.WriteString(w, `,`)
			stdprinter.PanicIfError(err)
		}

		item := alfredItemFromTag(t, c)
		j, err := lib.GenerateAlfredItemsJson([]lib.AlfredItem{item})

		if err != nil {
			//TODO-P2 handle error
			continue
		}

		_, err = io.WriteString(w, j)
		stdprinter.PanicIfError(err)

		needsComma = true
	}

	_, err = io.WriteString(w, `]}`)
	stdprinter.PanicIfError(err)
}
