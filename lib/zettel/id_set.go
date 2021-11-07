package zettel

import (
	"encoding/json"
	"sync"

	"github.com/friedenberg/z/util/stdprinter"
)

func MakeIdSet() *IdSet {
	return &IdSet{
		Locker: &sync.Mutex{},
		Set:    make(map[Id]bool),
	}
}

type IdSet struct {
	sync.Locker `json:"-"`
	Set         map[Id]bool
}

func (s *IdSet) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &s.Set); err != nil {
		return err
	}

	s.Locker = &sync.Mutex{}

	return nil
}

func (s *IdSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Set)
}

func (s *IdSet) Add(in ...Id) {
	s.Lock()
	defer s.Unlock()

	for _, i := range in {
		stdprinter.Debug("adding id to set:", i)
		s.Set[i] = true
	}

	return
}

func (s *IdSet) Delete(i Id) {
	s.Lock()
	defer s.Unlock()

	stdprinter.Debug("deleting id from set:", i)
	delete(s.Set, i)
}

func (s *IdSet) Len() int {
	return len(s.Set)
}

func (s *IdSet) Slice() (out []Id) {
	s.Lock()
	defer s.Unlock()

	for k, _ := range s.Set {
		out = append(out, k)
	}

	return
}
