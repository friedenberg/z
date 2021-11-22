package zettel

import (
	"reflect"
	"sync"
	"testing"
)

func TestMultiMapAdd(t *testing.T) {
	l := &sync.Mutex{}
	m := MakeMultiMap(l)

	m.Add("wow", Id(0))
	expected := []string{"wow"}
	assertValuesEqual(t, m, Id(0), expected)
}

func TestMultiMapAddThenDelete(t *testing.T) {
	l := &sync.Mutex{}
	m := MakeMultiMap(l)

	m.Add("wow", Id(0))

	assertValuesEqual(t, m, Id(0), []string{"wow"})
	assertIdsEqual(t, m, "wow", []Id{Id(0)})

	m.Delete(Id(0))

	assertValuesEmpty(t, m, Id(0))
	assertIdsEmpty(t, m, "wow")
}

func TestMultiMapAddSeveralThenDelete(t *testing.T) {
	l := &sync.Mutex{}
	m := MakeMultiMap(l)

	m.Add("wow", Id(0))
	m.Add("wow", Id(1))

	assertValuesEqual(t, m, Id(0), []string{"wow"})
	assertValuesEqual(t, m, Id(1), []string{"wow"})
	assertIdsEqual(t, m, "wow", []Id{Id(0), Id(1)})

	m.Delete(Id(0))

	assertValuesEmpty(t, m, Id(0))
	assertValuesEqual(t, m, Id(1), []string{"wow"})
	assertIdsEqual(t, m, "wow", []Id{Id(1)})
}

func assertIdsEmpty(t *testing.T, m MultiMap, k string) {
	t.Helper()
	actual1, ok := m.GetIds(k)

	if ok {
		t.Errorf("expected no values, but got some values: %#v", actual1.Set)
	}
}

func assertValuesEmpty(t *testing.T, m MultiMap, id Id) {
	t.Helper()
	actual1, ok := m.GetValues(id)

	if ok {
		t.Errorf("expected no values, but got some values: %#v", actual1)
	}
}

func assertValuesEqual(t *testing.T, m MultiMap, id Id, expected []string) {
	t.Helper()
	actual1, ok := m.GetValues(id)

	if !ok {
		t.Errorf("expected some values, but got none at all")
		return
	}

	actual := actual1.Slice()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %s but got %s", expected, actual)
	}
}

func assertIdsEqual(t *testing.T, m MultiMap, k string, expected []Id) {
	t.Helper()
	actual1, ok := m.GetIds(k)

	if !ok {
		t.Errorf("expected some values, but got none at all")
		return
	}

	actual := actual1.Slice()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %s but got %s", expected, actual)
	}
}
