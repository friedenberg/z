package lib

import (
	"encoding/json"
	"io"
	"sync"
	"time"

	"github.com/friedenberg/z/lib/collections"
	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/metadata"
	"golang.org/x/xerrors"
)

type IndexZettel struct {
	Path     string
	Id       int64
	Metadata metadata.Metadata
	//TODO-P2 remove
	Body string
}

//TODO-P4 move into index and use unmarshal methods
type SerializableIndex struct {
	ModTime int64
	Zettels map[zettel.Id]IndexZettel
	Files   collections.Map
	Urls    collections.Map
	Tags    collections.MultiMap
}

type Index struct {
	SerializableIndex
	*sync.Mutex
}

func MakeIndex() *Index {
	return &Index{
		SerializableIndex: SerializableIndex{
			Zettels: make(map[zettel.Id]IndexZettel),
			Files:   collections.MakeMap(),
			Urls:    collections.MakeMap(),
			Tags:    collections.MakeMultiMap(),
		},
		Mutex: &sync.Mutex{},
	}
}

func (i Index) Read(r io.Reader) (err error) {
	dec := json.NewDecoder(r)
	// dec := gob.NewDecoder(r)
	err = dec.Decode(&i.SerializableIndex)

	if err != nil {
		return
	}

	return
}

func (i Index) Write(w io.Writer) (err error) {
	i.ModTime = time.Now().Unix()
	enc := json.NewEncoder(w)
	// enc := gob.NewEncoder(w)
	err = enc.Encode(i.SerializableIndex)

	if err != nil {
		return
	}

	return
}

func (m Index) Get(k zettel.Id) (IndexZettel, bool) {
	if m.Mutex == nil {
		panic("mutex was not initalized")
	}

	m.Lock()
	defer m.Unlock()
	a, ok := m.Zettels[k]
	return a, ok
}

func (m Index) set(k zettel.Id, z IndexZettel) {
	m.Lock()
	defer m.Unlock()
	m.Zettels[k] = z
}

//TODO-P0 check for checksum file name collisions
func (i Index) Add(z *Zettel) error {
	if _, ok := i.Get(zettel.Id(z.Id)); ok {
		return xerrors.Errorf("zettel with id '%d' already exists in index", z.Id)
	}

	i.set(zettel.Id(z.Id), IndexZettel{
		Path:     z.Path,
		Id:       int64(z.Id),
		Metadata: z.Metadata,
		Body:     z.Body,
	})

	if u, ok := z.Metadata.Url(); ok {
		//TODO-P3 ensure not clobbering
		i.Urls.Set(u.String(), zettel.Id(z.Id), i)
	}

	for _, t := range z.Metadata.TagStrings() {
		i.Tags.Add(t, zettel.Id(z.Id), i)
	}

	return nil
}

func (i Index) AddFile(z *Zettel, sum string) (err error) {
	i.Files.Set(sum, zettel.Id(z.Id), i)

	return
}

func (i Index) Update(z *Zettel) (err error) {
	err = i.Delete(z)

	if err != nil {
		return err
	}

	err = i.Add(z)

	return
}

func (i Index) Delete(z *Zettel) (err error) {
	id := zettel.Id(z.Id)
	delete(i.Zettels, id)
	i.Files.Delete(id, i)
	i.Urls.Delete(id, i)
	i.Tags.Delete(id, i)
	return
}

func (i Index) HydrateZettel(z *Zettel, zb IndexZettel) {
	z.Metadata = zb.Metadata
	z.Id = zb.Id
	z.Path = zb.Path
	z.Body = zb.Body
}
