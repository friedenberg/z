package collections

import (
	"encoding/json"
	"sync"

	"github.com/friedenberg/z/lib/zettel"
)

func MakeMultiMap(l sync.Locker) MultiMap {
	return MultiMap{
		Locker: l,
		multiMapSerializable: multiMapSerializable{

			ValueToId: make(map[string]*ZettelIdSet),
			IdToValue: make(map[zettel.Id]*StringSet),
		},
	}
}

type multiMapSerializable struct {
	ValueToId map[string]*ZettelIdSet
	IdToValue map[zettel.Id]*StringSet
}

type MultiMap struct {
	sync.Locker
	multiMapSerializable
}

func (s *MultiMap) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &s.multiMapSerializable); err != nil {
		return err
	}

	return nil
}

func (s MultiMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.multiMapSerializable)
}

func (m MultiMap) GetIds(k string) (*ZettelIdSet, bool) {
	m.Lock()
	defer m.Unlock()

	a, ok := m.ValueToId[k]
	return a, ok
}

func (m MultiMap) GetValues(id zettel.Id) (*StringSet, bool) {
	m.Lock()
	defer m.Unlock()

	a, ok := m.IdToValue[id]
	return a, ok
}

func (m MultiMap) Set(k string, ids *ZettelIdSet) {
	m.Lock()
	defer m.Unlock()

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

func (m MultiMap) Add(k string, id zettel.Id) {
	a, _ := m.GetIds(k)

	if a == nil {
		a = MakeZettelIdSet()
	}

	a.Add(id)
	m.Set(k, a)
}

func (m MultiMap) Delete(id zettel.Id) {
	m.Lock()
	defer m.Unlock()

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
