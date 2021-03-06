package lib

import (
	"encoding/json"
	"io"
	"sync"
	"time"

	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/metadata"
	"github.com/friedenberg/z/util/stdprinter"
)

type IndexZettel struct {
	Path     string
	Id       zettel.Id
	Metadata metadata.Metadata
	Body     string
}

//TODO-P4 move into index and use unmarshal methods
type serializableIndex struct {
	ModTime int64
	Zettels map[zettel.Id]IndexZettel
	Files   zettel.Map
	Urls    zettel.Map
	Tags    zettel.MultiMap
}

type Index struct {
	serializableIndex
	*sync.Mutex
}

func MakeIndex() *Index {
	m := &sync.Mutex{}

	return &Index{
		serializableIndex: serializableIndex{
			Zettels: make(map[zettel.Id]IndexZettel),
			Files:   zettel.MakeMap(m),
			Urls:    zettel.MakeMap(m),
			Tags:    zettel.MakeMultiMap(m),
		},
		Mutex: m,
	}
}

func (i Index) Read(r io.Reader) (err error) {
	dec := json.NewDecoder(r)
	// dec := gob.NewDecoder(r)
	err = dec.Decode(&i.serializableIndex)

	if err != nil {
		return
	}

	i.Files.Locker = i
	i.Urls.Locker = i
	i.Tags.Locker = i

	return
}

func (i Index) Write(w io.Writer) (err error) {
	stdprinter.Debug("will write index")
	i.ModTime = time.Now().Unix()
	enc := json.NewEncoder(w)
	// enc := gob.NewEncoder(w)
	err = enc.Encode(i.serializableIndex)

	if err != nil {
		return
	}

	stdprinter.Debug("did write index")
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
func (i Index) Add(z *zettel.Zettel) {
	stdprinter.Debug("will add zettel in index:", z.Path)
	if _, ok := i.Get(zettel.Id(z.Id)); ok {
		stdprinter.Debugf("zettel with id '%d' already exists in index\n", z.Id)
	}

	i.set(zettel.Id(z.Id), IndexZettel{
		Path:     z.Path,
		Id:       z.Id,
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

	stdprinter.Debug("did add zettel in index:", z.Path)
}

func (i Index) AddFile(z *zettel.Zettel, sum string) (err error) {
	stdprinter.Debug("will add sum to index:", z.Path, sum)
	i.Files.Set(sum, zettel.Id(z.Id))
	stdprinter.Debug("did add sum to index:", z.Path, sum)

	return
}

func (i Index) Update(z *zettel.Zettel) (err error) {
	stdprinter.Debug("will update zettel in index:", z.Path)
	err = i.Delete(z)

	if err != nil {
		return err
	}

	i.Add(z)

	stdprinter.Debug("did update zettel in index:", z.Path)
	return
}

func (i Index) Delete(z *zettel.Zettel) (err error) {
	stdprinter.Debug("will delete zettel in index:", z.Path)
	id := zettel.Id(z.Id)
	delete(i.Zettels, id)
	i.Files.Delete(id)
	i.Urls.Delete(id)
	i.Tags.Delete(id)
	stdprinter.Debug("did delete zettel in index:", z.Path)
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
