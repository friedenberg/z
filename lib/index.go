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
	m := &sync.Mutex{}

	return &Index{
		SerializableIndex: SerializableIndex{
			Zettels: make(map[zettel.Id]IndexZettel),
			Files:   collections.MakeMap(m),
			Urls:    collections.MakeMap(m),
			Tags:    collections.MakeMultiMap(m),
		},
		Mutex: m,
	}
}

func (i Index) Read(r io.Reader) (err error) {
	dec := json.NewDecoder(r)
	// dec := gob.NewDecoder(r)
	err = dec.Decode(&i.SerializableIndex)

	if err != nil {
		return
	}

	i.Files.Locker = i
	i.Urls.Locker = i
	i.Tags.Locker = i

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
func (i Index) Add(z *zettel.Zettel) error {
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
		i.Urls.Set(u.String(), zettel.Id(z.Id))
	}

	for _, t := range z.Metadata.TagStrings() {
		i.Tags.Add(t, zettel.Id(z.Id))
	}

	return nil
}

func (i Index) AddFile(z *zettel.Zettel, sum string) (err error) {
	i.Files.Set(sum, zettel.Id(z.Id))

	return
}

func (i Index) Update(z *zettel.Zettel) (err error) {
	err = i.Delete(z)

	if err != nil {
		return err
	}

	err = i.Add(z)

	return
}

func (i Index) Delete(z *zettel.Zettel) (err error) {
	id := zettel.Id(z.Id)
	delete(i.Zettels, id)
	i.Files.Delete(id)
	i.Urls.Delete(id)
	i.Tags.Delete(id)
	return
}

func (i Index) ForFileSum(sum string) (z *zettel.Zettel, ok bool) {
	oldZettelId, ok := i.Files.GetId(sum)

	if ok {
		iz, ok := i.Get(oldZettelId)

		if !ok {
			panic("index had file in zettel but no zettel in index")
		}

		z = &zettel.Zettel{}
		i.HydrateZettel(z, iz)
	}

	return
}

func (i Index) HydrateZettel(z *zettel.Zettel, zb IndexZettel) {
	z.Metadata = zb.Metadata
	z.Id = zb.Id
	z.Path = zb.Path
	z.Body = zb.Body
}
