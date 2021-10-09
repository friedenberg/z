package lib

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/friedenberg/z/lib/zettel"
	"golang.org/x/xerrors"
)

type IndexZettel struct {
	Path     string
	Id       int64
	Metadata Metadata
	//TODO remove
	Body string
}

type SerializableIndex struct {
	Zettels map[zettel.Id]IndexZettel
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
			Zettels: make(map[zettel.Id]IndexZettel),
			Files:   MakeZettelIdMap(),
			Urls:    MakeZettelIdMap(),
			Tags:    MakeZettelIdMap(),
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

	if z.HasFile() {
		i.Files.Add(z.FilePath(), zettel.Id(z.Id), i)
	}

	if z.HasUrl() {
		i.Urls.Add(z.Metadata.Url, zettel.Id(z.Id), i)
	}

	for _, t := range z.Metadata.Tags {
		i.Tags.Add(t, zettel.Id(z.Id), i)
	}

	return nil
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

func (i Index) ZettelsForUrl(u string) (o []IndexZettel) {
	//TODO normalize url
	ids, ok := i.Urls.Get(u, i)

	if !ok {
		return
	}

	for _, id := range ids.Slice() {
		if zi, ok := i.Zettels[zettel.Id(id)]; ok {
			o = append(o, zi)
		}
	}

	return
}
