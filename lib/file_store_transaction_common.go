package lib

import "os"

func (k *FileStore) readAndWrite(z *Zettel, shouldWrite bool) (err error) {
	err = k.Hydrate(z, true)

	if os.IsNotExist(err) {
		err = nil
	} else if err != nil {
		return
	}

	// if f, ok := z.Note.Metadata.NewFile(); ok {
	// 	var sum string
	// 	sum, err = util.Sha256HashForFile(f.Path)

	// 	if err != nil {
	// 		return
	// 	}

	// 	z, ok := u.Index.GetZettelForFileSha(sum)

	// 	if ok {
	// 	}

	// 	fd := metadata.File{
	// 		Id:  sum,
	// 		Ext: f.Ext(),
	// 	}

	// 	if t.ShouldCopyFiles {
	// 		cmd := exec.Command("cp", "-R", file, fd.FilePath(u.BasePath))
	// 		msg, err := cmd.CombinedOutput()

	// 		if err != nil {
	// 			err = xerrors.Errorf("%w: %s", err, msg)
	// 		}
	// 	} else {
	// 		cmd := exec.Command("mv", file, fd.FilePath(u.BasePath))
	// 		msg, err := cmd.CombinedOutput()

	// 		if err != nil {
	// 			err = xerrors.Errorf("%w: %s", err, msg)
	// 		}
	// 	}

	// 	if err != nil {
	// 		return
	// 	}
	// }

	if !shouldWrite {
		return
	}

	err = z.Write(nil)

	if err != nil {
		return
	}

	return
}
