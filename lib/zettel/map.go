package zettel

import (
	"encoding/json"
	"sync"
)

func MakeMap(l sync.Locker) Map {
	return Map{
		Locker: l,
		mapSerializable: mapSerializable{
			ValueToId: make(map[string]Id),
			IdToValue: make(map[Id]string),
		},
	}
}

type mapSerializable struct {
	ValueToId map[string]Id
	IdToValue map[Id]string
}

type Map struct {
	sync.Locker
	mapSerializable
}

func (s *Map) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &s.mapSerializable); err != nil {
		return err
	}

	return nil
}

func (s Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.mapSerializable)
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

	delete(m.IdToValue, id)
	delete(m.ValueToId, v)
}
