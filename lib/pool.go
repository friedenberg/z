package lib

import (
	"sync"
)

type ZettelPool interface {
	Get() *Zettel
	Put(*Zettel)
}

type zettelPool struct {
	sync.Pool
	env *Env
}

func (p zettelPool) Get() (z *Zettel) {
	z = p.Pool.Get().(*Zettel)
	z.IndexData.Date = ""
	z.IndexData.Description = ""
	z.IndexData.Tags = z.IndexData.Tags[:0]
	z.IndexData.Url = ""
	// z.AlfredData.Item.Title = ""
	// z.AlfredData.Item.Arg = ""
	// z.AlfredData.Item.Subtitle = ""
	// z.AlfredData.Item.Match = ""
	// z.AlfredData.Item.Icon = alfred.AlfredItemIcon{}
	// z.AlfredData.Item.Uid = ""
	// z.AlfredData.Item.ItemType = ""
	// z.AlfredData.Item.QuicklookUrl = ""
	// z.AlfredData.Item.Text.Copy = ""
	return
}

func (p zettelPool) Put(z *Zettel) {
	p.Pool.Put(z)
}
