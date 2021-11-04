package lib

import (
	"sync"

	"github.com/friedenberg/z/lib/zettel"
)

type mapZettelIdToZettel map[zettel.Id]*Zettel

func MakeSet() Set {
	return Set{
		Locker: &sync.Mutex{},
		set:    make(mapZettelIdToZettel),
	}
}

type Set struct {
	sync.Locker
	set mapZettelIdToZettel
}

func (m Set) Get(i zettel.Id) (*Zettel, bool) {
	m.Lock()
	defer m.Unlock()

	a, ok := m.set[i]
	return a, ok
}

func (m Set) ModifyZettel(i int, z *Zettel) (err error) {
	m.Add(z)
	return
}

func (m Set) Add(z *Zettel) {
	m.Lock()
	defer m.Unlock()

	m.set[zettel.Id(z.Id)] = z
}

func (m Set) Del(z *Zettel) {
	m.Lock()
	defer m.Unlock()

	delete(m.set, zettel.Id(z.Id))
}

func (m Set) Len() int {
	return len(m.set)
}

func (m Set) Zettels() (zs ZettelSlice) {
	zs = make(ZettelSlice, 0, len(m.set))

	for _, z := range m.set {
		zs = append(zs, z)
	}

	return
}
