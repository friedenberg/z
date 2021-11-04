package collections

import (
	"encoding/json"
	"sync"

	"github.com/friedenberg/z/lib/zettel"
)

func MakeMap(l sync.Locker) Map {
	return Map{
		Locker: l,
		mapSerializable: mapSerializable{
			ValueToId: make(map[string]zettel.Id),
			IdToValue: make(map[zettel.Id]string),
		},
	}
}

type mapSerializable struct {
	ValueToId map[string]zettel.Id
	IdToValue map[zettel.Id]string
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

func (m Map) GetId(k string) (zettel.Id, bool) {
	m.Lock()
	defer m.Unlock()
	a, ok := m.ValueToId[k]
	return a, ok
}

func (m Map) GetValue(id zettel.Id) (string, bool) {
	m.Lock()
	defer m.Unlock()
	a, ok := m.IdToValue[id]
	return a, ok
}

func (m Map) Set(k string, id zettel.Id) {
	m.Lock()
	defer m.Unlock()

	m.ValueToId[k] = id
	m.IdToValue[id] = k
}

func (m Map) Delete(id zettel.Id) {
	m.Lock()
	defer m.Unlock()
	v, ok := m.IdToValue[id]

	if !ok {
		return
	}

	delete(m.IdToValue, id)
	delete(m.ValueToId, v)
}
