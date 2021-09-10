package lib

import (
	"encoding/gob"
	"fmt"
	"io"
	"strconv"
	"sync"
)

var thisFileSha1 string

type IndexZettel struct {
	Path     string
	Id       int64
	Metadata Metadata
}

type ZettelIdMap struct {
	IdMap map[string][]int64
}

func MakeZettelIdMap() ZettelIdMap {
	return ZettelIdMap{
		IdMap: make(map[string][]int64),
	}
}

func (m ZettelIdMap) Get(k string, l sync.Locker) ([]int64, bool) {
	l.Lock()
	defer l.Unlock()
	a, ok := m.IdMap[k]
	return a, ok
}

func (m ZettelIdMap) Set(k string, ids []int64, l sync.Locker) {
	l.Lock()
	defer l.Unlock()
	m.IdMap[k] = ids
}

func (m ZettelIdMap) Add(k string, id int64, l sync.Locker) {
	var a []int64
	ok := false

	if a, ok = m.Get(k, l); !ok {
		a = make([]int64, 0, 1)
	}

	a = append(a, id)
	m.Set(k, a, l)
}

type SerializableIndex struct {
	Version string
	Zettels map[string]IndexZettel
	Files   ZettelIdMap
	Urls    ZettelIdMap
	Tags    ZettelIdMap
}

type Index struct {
	SerializableIndex
	*sync.Mutex
}

func MakeIndex() Index {
	return Index{
		SerializableIndex: SerializableIndex{
			Version: thisFileSha1,
			Zettels: make(map[string]IndexZettel),
			Files:   MakeZettelIdMap(),
			Urls:    MakeZettelIdMap(),
			Tags:    MakeZettelIdMap(),
		},
		Mutex: &sync.Mutex{},
	}
}

func (i Index) Read(r io.Reader) (err error) {
	dec := gob.NewDecoder(r)
	err = dec.Decode(&i.SerializableIndex)

	if err != nil {
		return
	}

	return
}

func (i Index) Write(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	err = enc.Encode(i.SerializableIndex)

	if err != nil {
		return
	}

	return
}

func (m Index) Get(k string) (IndexZettel, bool) {
	m.Lock()
	defer m.Unlock()
	a, ok := m.Zettels[k]
	return a, ok
}

func (m Index) Set(k string, z IndexZettel) {
	m.Lock()
	defer m.Unlock()
	m.Zettels[k] = z
}

func (i Index) Add(z *Zettel) error {
	if _, ok := i.Get(strconv.FormatInt(z.Id, 10)); ok {
		return fmt.Errorf("zettel with id '%d' already exists in index", z.Id)
	}

	i.Set(strconv.FormatInt(z.Id, 10), IndexZettel{
		Path:     z.Path,
		Id:       z.Id,
		Metadata: z.Metadata,
	})

	if z.HasFile() {
		i.Files.Add(z.FilePath(), z.Id, i)
	}

	if z.HasUrl() {
		i.Urls.Add(z.Metadata.Url, z.Id, i)
	}

	for _, t := range z.Metadata.Tags {
		i.Tags.Add(t, z.Id, i)
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
	ids, ok := i.Urls.Get(u, i)

	if !ok {
		return
	}

	for _, id := range ids {
		if zi, ok := i.Zettels[strconv.FormatInt(id, 10)]; ok {
			o = append(o, zi)
		}
	}

	return
}
