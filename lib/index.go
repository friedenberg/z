package lib

import (
	"fmt"
	"sync"
)

type IndexZettel struct {
	Path     string
	Id       int64
	Metadata Metadata
}

type ZettelIdMap struct {
	idMap map[string][]int64
	*sync.Mutex
}

func MakeZettelIdMap() ZettelIdMap {
	return ZettelIdMap{
		idMap: make(map[string][]int64),
		Mutex: &sync.Mutex{},
	}
}

func (m ZettelIdMap) Get(k string) ([]int64, bool) {
	m.Lock()
	defer m.Unlock()
	a, ok := m.idMap[k]
	return a, ok
}

func (m ZettelIdMap) Set(k string, ids []int64) {
	m.Lock()
	defer m.Unlock()
	m.idMap[k] = ids
}

func (m ZettelIdMap) Add(k string, id int64) {
	var a []int64
	ok := false

	if a, ok = m.Get(k); !ok {
		a = make([]int64, 0, 1)
	}

	a = append(a, id)
	m.Set(k, a)
}

type Index struct {
	Zettels map[int64]IndexZettel
	Files   ZettelIdMap
	Urls    ZettelIdMap
	Tags    ZettelIdMap
	*sync.Mutex
}

func MakeIndex() Index {
	return Index{
		Zettels: make(map[int64]IndexZettel),
		Files:   MakeZettelIdMap(),
		Urls:    MakeZettelIdMap(),
		Tags:    MakeZettelIdMap(),
		Mutex:   &sync.Mutex{},
	}
}

func (m Index) Get(k int64) (IndexZettel, bool) {
	m.Lock()
	defer m.Unlock()
	a, ok := m.Zettels[k]
	return a, ok
}

func (m Index) Set(k int64, z IndexZettel) {
	m.Lock()
	defer m.Unlock()
	m.Zettels[k] = z
}

func (i Index) Add(z *Zettel) error {
	if _, ok := i.Get(z.Id); ok {
		return fmt.Errorf("zettel with id '%d' already exists in index", z.Id)
	}

	i.Set(z.Id, IndexZettel{
		Path:     z.Path,
		Id:       z.Id,
		Metadata: z.Metadata,
	})

	if z.HasFile() {
		i.Files.Add(z.FilePath(), z.Id)
	}

	if z.HasUrl() {
		i.Urls.Add(z.Metadata.Url, z.Id)
	}

	for _, t := range z.Metadata.Tags {
		i.Tags.Add(t, z.Id)
	}

	return nil
}

func (i Index) HydrateZettel(z *Zettel, zb IndexZettel) {
	z.Metadata = zb.Metadata
	z.Id = zb.Id
	z.Path = zb.Path
}

func (i Index) ZettelsForUrl(u string) (o []IndexZettel) {
	//TODO normalize url
	ids, ok := i.Urls.Get(u)

	if !ok {
		return
	}

	for _, id := range ids {
		if zi, ok := i.Zettels[id]; ok {
			o = append(o, zi)
		}
	}

	return
}
