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
	z.IndexData.Date = ""
	z.IndexData.Description = ""
	z.IndexData.Areas = z.IndexData.Areas[:0]
	z.IndexData.Projects = z.IndexData.Projects[:0]
	z.IndexData.Tags = z.IndexData.Tags[:0]
	z.IndexData.Url = ""
	z.AlfredData.Item.Title = ""
	z.AlfredData.Item.Arg = ""
	z.AlfredData.Item.Subtitle = ""
	z.AlfredData.Item.Match = ""
	z.AlfredData.Item.Icon = ZettelAlfredItemIcon{}
	z.AlfredData.Item.Uid = ""
	z.AlfredData.Item.ItemType = ""
	z.AlfredData.Item.QuicklookUrl = ""
	z.AlfredData.Item.Text.Copy = ""
	return
}
