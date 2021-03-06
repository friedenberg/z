package metadata

import (
	"sort"
)

func MakeTagSet() TagSet {
	return TagSet{
		set: make(map[string]ITag),
	}
}

type TagSet struct {
	set map[string]ITag
}

func (s TagSet) addOne(t ITag) (ok bool) {
	ts := t.Tag()
	old, ok := s.set[ts]

	if ok && old != t {
		ok = false
		return
	}

	s.set[ts] = t

	return
}

func (s TagSet) Merge(s1 TagSet) {
	for _, t := range s1.Tags() {
		s.addOne(t)
	}

	return
}

func (s TagSet) Add(t ...ITag) (ok bool) {
	for _, t1 := range t {
		ok = s.addOne(t1)

		if !ok {
			return
		}
	}

	return
}

func (s TagSet) Get(q string) (t1 ITag, ok bool) {
	t1, ok = s.set[q]
	return
}

func (s TagSet) Del(t string) {
	delete(s.set, t)
	return
}

func (s TagSet) Len() int {
	return len(s.set)
}

func (s TagSet) Tags() (ts []ITag) {
	ts = make([]ITag, 0, len(s.set))

	for _, t := range s.set {
		ts = append(ts, t)
	}

	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Tag() < ts[j].Tag()
	})

	return
}

func (s TagSet) Strings() (ts []string) {
	ts = make([]string, 0, len(s.set))

	for t, _ := range s.set {
		ts = append(ts, t)
	}

	sort.Slice(ts, func(i, j int) bool {
		return ts[i] < ts[j]
	})

	return
}
