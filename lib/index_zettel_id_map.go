package lib

import (
	"sync"

	"github.com/friedenberg/z/lib/set"
	"github.com/friedenberg/z/lib/zettel"
)

func MakeZettelIdMap() ZettelIdMap {
	return ZettelIdMap{
		ValueToId: make(map[string]*set.ZettelIdSet),
		IdToValue: make(map[zettel.Id]*set.StringSet),
	}
}

type ZettelIdMap struct {
	ValueToId map[string]*set.ZettelIdSet
	IdToValue map[zettel.Id]*set.StringSet
}

func (m ZettelIdMap) Get(k string, l sync.Locker) (*set.ZettelIdSet, bool) {
	l.Lock()
	defer l.Unlock()
	a, ok := m.ValueToId[k]
	return a, ok
}

func (m ZettelIdMap) Set(k string, ids *set.ZettelIdSet, l sync.Locker) {
	l.Lock()
	defer l.Unlock()

	m.ValueToId[k] = ids

	for _, id := range ids.Slice() {
		a, _ := m.IdToValue[id]

		if a == nil {
			a = set.MakeStringSet()
		}

		a.Add(k)
		m.IdToValue[id] = a
	}
}

func (m ZettelIdMap) Add(k string, id zettel.Id, l sync.Locker) {
	a, _ := m.Get(k, l)

	if a == nil {
		a = set.MakeZettelIdSet()
	}

	a.Add(id)
	m.Set(k, a, l)
}

func (m ZettelIdMap) Delete(id zettel.Id, l sync.Locker) {
	l.Lock()
	defer l.Unlock()
	vs, ok := m.IdToValue[id]

	if !ok {
		return
	}

	delete(m.IdToValue, id)

	if vs == nil {
		//TODO determine why this is necessary
		return
	}

	for _, v := range vs.Slice() {
		delete(m.ValueToId, v)
	}
}
