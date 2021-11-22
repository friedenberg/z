package zettel

import (
	"encoding/json"
	"sync"

	"github.com/friedenberg/z/util/stdprinter"
)

func MakeMap(l sync.Locker) Map {
	return Map{
		Locker: l,
		MapSerializable: MapSerializable{
			ValueToId: make(map[string]Id),
			IdToValue: make(map[Id]string),
		},
	}
}

type MapSerializable struct {
	ValueToId map[string]Id
	IdToValue map[Id]string
}

type Map struct {
	sync.Locker
	MapSerializable
}

func (s *Map) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &s.MapSerializable); err != nil {
		return err
	}

	return nil
}

func (s Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.MapSerializable)
}

func (m Map) GetId(k string) (Id, bool) {
	m.Lock()
	defer m.Unlock()
	a, ok := m.ValueToId[k]
	return a, ok
}

func (m Map) GetValue(id Id) (string, bool) {
	m.Lock()
	defer m.Unlock()
	a, ok := m.IdToValue[id]
	return a, ok
}

func (m Map) Set(k string, id Id) {
	m.Lock()
	defer m.Unlock()

	stdprinter.Debug("setting key and id in map:", k, id)
	m.ValueToId[k] = id
	m.IdToValue[id] = k
}

func (m Map) Delete(id Id) {
	m.Lock()
	defer m.Unlock()
	v, ok := m.IdToValue[id]

	if !ok {
		return
	}

	stdprinter.Debug("deleting id from map:", id)
	delete(m.IdToValue, id)
	delete(m.ValueToId, v)
}
