package collections

import (
	"sync"

	"github.com/friedenberg/z/lib/zettel"
)

func MakeMap() Map {
	return Map{
		ValueToId: make(map[string]zettel.Id),
		IdToValue: make(map[zettel.Id]string),
	}
}

//TODO add option for limiting set to one value
type Map struct {
	ValueToId map[string]zettel.Id
	IdToValue map[zettel.Id]string
}

func (m Map) GetId(k string, l sync.Locker) (zettel.Id, bool) {
	l.Lock()
	defer l.Unlock()
	a, ok := m.ValueToId[k]
	return a, ok
}

func (m Map) GetValue(id zettel.Id, l sync.Locker) (string, bool) {
	l.Lock()
	defer l.Unlock()
	a, ok := m.IdToValue[id]
	return a, ok
}

func (m Map) Set(k string, id zettel.Id, l sync.Locker) {
	l.Lock()
	defer l.Unlock()

	m.ValueToId[k] = id
	m.IdToValue[id] = k
}

func (m Map) Delete(id zettel.Id, l sync.Locker) {
	l.Lock()
	defer l.Unlock()
	v, ok := m.IdToValue[id]

	if !ok {
		return
	}

	delete(m.IdToValue, id)
	delete(m.ValueToId, v)
}
