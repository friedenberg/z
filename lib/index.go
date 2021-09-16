package lib

import (
	"encoding/gob"
	"io"
	"sync"

	"golang.org/x/xerrors"
)

type IndexZettel struct {
	Path     string
	Id       int64
	Metadata Metadata
	Body     string
}

type ZettelIdMap struct {
	IdMap map[string][]Id
}

func MakeZettelIdMap() ZettelIdMap {
	return ZettelIdMap{
		IdMap: make(map[string][]Id),
	}
}

func (m ZettelIdMap) Get(k string, l sync.Locker) ([]Id, bool) {
	l.Lock()
	defer l.Unlock()
	a, ok := m.IdMap[k]
	return a, ok
}

func (m ZettelIdMap) Set(k string, ids []Id, l sync.Locker) {
	l.Lock()
	defer l.Unlock()
	m.IdMap[k] = ids
}

func (m ZettelIdMap) Add(k string, id Id, l sync.Locker) {
	var a []Id
	ok := false

	if a, ok = m.Get(k, l); !ok {
		a = make([]Id, 0, 1)
	}

	a = append(a, id)
	m.Set(k, a, l)
}

type SerializableIndex struct {
	Zettels map[Id]IndexZettel
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
			Zettels: make(map[Id]IndexZettel),
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

func (m Index) Get(k Id) (IndexZettel, bool) {
	m.Lock()
	defer m.Unlock()
	a, ok := m.Zettels[k]
	return a, ok
}

func (m Index) Set(k Id, z IndexZettel) {
	m.Lock()
	defer m.Unlock()
	m.Zettels[k] = z
}

func (i Index) Add(z *Zettel) error {
	if _, ok := i.Get(Id(z.Id)); ok {
		return xerrors.Errorf("zettel with id '%d' already exists in index", z.Id)
	}

	i.Set(Id(z.Id), IndexZettel{
		Path:     z.Path,
		Id:       int64(z.Id),
		Metadata: z.Metadata,
		Body:     z.Body,
	})

	if z.HasFile() {
		i.Files.Add(z.FilePath(), Id(z.Id), i)
	}

	if z.HasUrl() {
		i.Urls.Add(z.Metadata.Url, Id(z.Id), i)
	}

	for _, t := range z.Metadata.Tags {
		i.Tags.Add(t, Id(z.Id), i)
	}

	return nil
}

func (i Index) HydrateZettel(z *Zettel, zb IndexZettel) {
	z.Metadata = zb.Metadata
	z.Id = zb.Id
	z.Path = zb.Path
	z.Body = zb.Body
}

func (i Index) ZettelsForUrl(u string) (o []IndexZettel) {
	//TODO normalize url
	ids, ok := i.Urls.Get(u, i)

	if !ok {
		return
	}

	for _, id := range ids {
		if zi, ok := i.Zettels[Id(id)]; ok {
			o = append(o, zi)
		}
	}

	return
}
