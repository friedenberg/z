package lib

import (
	"sync"
)

type ZettelPool struct {
	sync.Pool
}

var ZettelPoolInstance = ZettelPool{
	sync.Pool{
		New: func() interface{} {
			return new(Zettel)
		},
	},
}

func (p ZettelPool) Get() (z *Zettel) {
	z = p.Pool.Get().(*Zettel)
	z.Metadata.Date = ""
	z.Metadata.Kind = ""
	z.Metadata.Description = ""
	z.Metadata.Areas = z.Metadata.Areas[:0]
	z.Metadata.Projects = z.Metadata.Projects[:0]
	z.Metadata.Tags = z.Metadata.Tags[:0]
	z.Metadata.Url = ""
	z.AlfredItem.Title = ""
	z.AlfredItem.Arg = ""
	z.AlfredItem.Subtitle = ""
	z.AlfredItem.Match = ""
	z.AlfredItem.Icon = ZettelAlfredItemIcon{}
	z.AlfredItem.Uid = ""
	z.AlfredItem.ItemType = ""
	z.AlfredItem.QuicklookUrl = ""
	z.AlfredItem.Text.Copy = ""
	return
}
