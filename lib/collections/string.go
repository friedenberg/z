package collections

import (
	"encoding/json"
	"sync"
)

func MakeStringSet() *StringSet {
	return &StringSet{
		Locker: &sync.Mutex{},
		Set:    make(map[string]bool),
	}
}

type StringSet struct {
	sync.Locker `json:"-"`
	Set         map[string]bool
}

func (s *StringSet) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &s.Set); err != nil {
		return err
	}

	s.Locker = &sync.Mutex{}

	return nil
}

func (s *StringSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Set)
}

func (s *StringSet) Add(in ...string) {
	s.Lock()
	defer s.Unlock()

	for _, i := range in {
		s.Set[i] = true
	}

	return
}

func (s *StringSet) Delete(i string) {
	s.Lock()
	defer s.Unlock()

	delete(s.Set, i)
}

func (s *StringSet) Len() int {
	return len(s.Set)
}

func (s *StringSet) Slice() (out []string) {
	s.Lock()
	defer s.Unlock()

	for k, _ := range s.Set {
		out = append(out, k)
	}

	return
}
