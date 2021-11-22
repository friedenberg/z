package zettel

import (
	"encoding/json"
	"sync"
)

func MakeMultiMap(l sync.Locker) MultiMap {
	return MultiMap{
		Locker: l,
		MultiMapSerializable: MultiMapSerializable{
			ValueToId: make(map[string]*IdSet),
			IdToValue: make(map[Id]*StringSet),
		},
	}
}

type MultiMapSerializable struct {
	ValueToId map[string]*IdSet
	IdToValue map[Id]*StringSet
}

type MultiMap struct {
	sync.Locker
	MultiMapSerializable
}

func (s *MultiMap) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &s.MultiMapSerializable); err != nil {
		return err
	}

	return nil
}

func (s MultiMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.MultiMapSerializable)
}

func (m MultiMap) GetIds(k string) (*IdSet, bool) {
	m.Lock()
	defer m.Unlock()

	a, ok := m.ValueToId[k]
	return a, ok
}

func (m MultiMap) GetValues(id Id) (*StringSet, bool) {
	m.Lock()
	defer m.Unlock()

	a, ok := m.IdToValue[id]
	return a, ok
}

func (m MultiMap) Set(k string, ids *IdSet) {
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

func (m MultiMap) Add(k string, id Id) {
	a, _ := m.GetIds(k)

	if a == nil {
		a = MakeIdSet()
	}

	a.Add(id)
	m.Set(k, a)
}

func (m MultiMap) Delete(id Id) {
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
		vs1, ok := m.ValueToId[v]

		if !ok {
			//TODO P4 this is inconsistent now?
		}

		vs1.Delete(id)

		if vs1.Len() == 0 {
			delete(m.ValueToId, v)
		}
	}
}
