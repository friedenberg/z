package lib

import (
	"os/exec"

	"github.com/friedenberg/z/lib/zettel"
	"github.com/friedenberg/z/lib/zettel/metadata"
	"github.com/friedenberg/z/util"
	"github.com/friedenberg/z/util/stdprinter"
	"golang.org/x/xerrors"
)

func (k *FileStore) transactionProcessAdd(u Umwelt, z *Zettel) (err error) {
	//add can be called for existing zettels or new zettels
	//in the case of new, we need to create an id and populate it
	if z.Id == 0 {
		var id zettel.Id
		id, err = u.Kasten.NewId()

		if err != nil {
			return
		}

		z.Id = id.Int()
		z.Path = MakePathFromId(u.Kasten.BasePath(), id.String())
	}

	err = k.hydrateFromFileIfExists(z)

	if err != nil {
		return
	}

	err = k.updateFilesIfNecessary(z)

	if err != nil {
		return
	}

	err = k.writeToFile(z)

	if err != nil {
		return
	}

	return
}

func (k *FileStore) updateFilesIfNecessary(z *Zettel) (err error) {
	lf, hasLocalFile := z.Note.Metadata.LocalFile()
	nf, hasNewFile := z.Note.Metadata.NewFile()

	if hasNewFile && hasLocalFile {
		err = xerrors.Errorf("zettel has both local and new file")
		return
	}

	if hasLocalFile {
		return k.updateLocalFile(z, lf)
	} else if hasNewFile {
		return k.updateNewFile(z, nf)
	}

	return
}

func (k *FileStore) updateNewFile(z *Zettel, f *metadata.NewFile) (err error) {
	var sum string
	sum, err = util.Sha256HashForFile(f.Path)

	if err != nil {
		err = xerrors.Errorf("failed to get sum for zettel: %s: %w", z.Id, err)
		return
	}

	oz, ok := k.umwelt.Index.ForFileSum(sum)

	if ok {
		oz.Merge(z)
		k.umwelt.Transaction.Set(oz, TransactionActionModified)
		k.umwelt.Transaction.Set(z, TransactionActionDeleted)
		z = oz
	} else {
		var f1 metadata.LocalFile
		f1, err = k.moveFile(z, f, sum)

		if err != nil {
			return
		}

		path := f1.FilePath(k.basePath)

		err = util.SetDisallowUserChanges(path)

		if err != nil {
			return
		}

		k.umwelt.Index.AddFile(z, sum)
	}

	return
}

func (k *FileStore) updateLocalFile(z *Zettel, f *metadata.LocalFile) (err error) {
	fPath := f.FilePath(k.basePath)

	isDir, err := util.IsDir(fPath)

	if isDir || err != nil {
		return
	}

	//TODO use real umwelt passed to this function
	oldSum, ok := k.umwelt.Index.Files.GetValue(zettel.Id(z.Id))

	if ok {
		//TODO: merge zettel
	} else {
		//TODO: handle case, is this possible?
	}

	var sum string
	stdprinter.Debugf("summing %s\n", fPath)
	sum, err = util.Sha256HashForFile(fPath)

	if err != nil {
		return
	}

	path := fPath

	if oldSum != "" && oldSum != sum {
		var f1 metadata.LocalFile
		f1, err = k.moveFile(z, f, sum)

		if err != nil {
			return
		}

		path = f1.FilePath(k.basePath)
	}

	err = util.SetDisallowUserChanges(path)

	if err != nil {
		return
	}

	k.umwelt.Index.AddFile(z, sum)

	return
}

func (k *FileStore) moveFile(z *Zettel, f metadata.File, sum string) (f1 metadata.LocalFile, err error) {
	fPath := f.FilePath(k.basePath)
	f1 = k.uniqueFile(sum, f.Extension())

	f1Path := f1.FilePath(k.basePath)

	stdprinter.Debugf("moving file:\n%#v\n%#v\n", f, f1)
	err = z.Metadata.AddStringTags(f1.Tag())

	if err != nil {
		return
	}

	err = util.SetAllowUserChanges(fPath)

	if err != nil {
		return
	}

	cmd := exec.Command("mv", fPath, f1Path)
	var msg []byte
	msg, err = cmd.CombinedOutput()

	if err != nil {
		err = xerrors.Errorf("%w: %s", err, msg)
		return
	}

	return
}

//TODO-P3 handle concurrency
func (k *FileStore) uniqueFile(sum, ext string) (f metadata.LocalFile) {
	f.Ext = ext

	for i := 7; i < 256; i++ {
		f.Id = sum[0:i]

		if !util.FileExists(f.FilePath(k.basePath)) {
			return
		}
	}

	err := xerrors.Errorf("unable to create unique id for %s.%s", sum, ext)
	stdprinter.PanicIfError(err)
	return
}
