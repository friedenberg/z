package zettel

import (
	"encoding/json"
	"sync"
)

func MakeZettelIdSet() *ZettelIdSet {
	return &ZettelIdSet{
		Locker: &sync.Mutex{},
		Set:    make(map[Id]bool),
	}
}

type ZettelIdSet struct {
	sync.Locker `json:"-"`
	Set         map[Id]bool
}

func (s *ZettelIdSet) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &s.Set); err != nil {
		return err
	}

	s.Locker = &sync.Mutex{}

	return nil
}

func (s *ZettelIdSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Set)
}

func (s *ZettelIdSet) Add(in ...Id) {
	s.Lock()
	defer s.Unlock()

	for _, i := range in {
		s.Set[i] = true
	}

	return
}

func (s *ZettelIdSet) Delete(i Id) {
	s.Lock()
	defer s.Unlock()

	delete(s.Set, i)
}

func (s *ZettelIdSet) Len() int {
	return len(s.Set)
}

func (s *ZettelIdSet) Slice() (out []Id) {
	s.Lock()
	defer s.Unlock()

	for k, _ := range s.Set {
		out = append(out, k)
	}

	return
}
