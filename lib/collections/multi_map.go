package collections

import (
	"sync"

	"github.com/friedenberg/z/lib/zettel"
)

func MakeMultiMap() MultiMap {
	return MultiMap{
		ValueToId: make(map[string]*ZettelIdSet),
		IdToValue: make(map[zettel.Id]*StringSet),
	}
}

//TODO add option for limiting set to one value
type MultiMap struct {
	ValueToId map[string]*ZettelIdSet
	IdToValue map[zettel.Id]*StringSet
}

func (m MultiMap) GetIds(k string, l sync.Locker) (*ZettelIdSet, bool) {
	l.Lock()
	defer l.Unlock()
	a, ok := m.ValueToId[k]
	return a, ok
}

func (m MultiMap) GetValues(id zettel.Id, l sync.Locker) (*StringSet, bool) {
	l.Lock()
	defer l.Unlock()
	a, ok := m.IdToValue[id]
	return a, ok
}

func (m MultiMap) Set(k string, ids *ZettelIdSet, l sync.Locker) {
	l.Lock()
	defer l.Unlock()

	m.ValueToId[k] = ids

	for _, id := range ids.Slice() {
		a, _ := m.IdToValue[id]

		if a == nil {
			a = MakeStringSet()
		}

		a.Add(k)
		m.IdToValue[id] = a
	}
}

func (m MultiMap) Add(k string, id zettel.Id, l sync.Locker) {
	a, _ := m.GetIds(k, l)

	if a == nil {
		a = MakeZettelIdSet()
	}

	a.Add(id)
	m.Set(k, a, l)
}

func (m MultiMap) Delete(id zettel.Id, l sync.Locker) {
	l.Lock()
	defer l.Unlock()
	vs, ok := m.IdToValue[id]

	if !ok {
		return
	}

	delete(m.IdToValue, id)

	if vs == nil {
		//TODO-P4 determine why this is necessary
		return
	}

	for _, v := range vs.Slice() {
		delete(m.ValueToId, v)
	}
}
