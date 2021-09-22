package kasten

import (
	"os/exec"
	"path"

	"golang.org/x/xerrors"
)

type Files struct {
	BasePath string
}

func (k *Files) InitFromOptions(o map[string]interface{}) (err error) {
	k.BasePath = getStringOption(o, "path")

	if k.BasePath == "" {
		err = xerrors.Errorf("config for 'files' kasten missing path")
		return
	}

	return
}

func (k *Files) CopyFileTo(localPath string, h RemoteFileHandle) (err error) {
	remotePath := path.Join(k.BasePath, h.FileName())

	cmd := exec.Command("cp", "-R", localPath, remotePath)
	out, err := cmd.CombinedOutput()

	if err != nil {
		err = xerrors.Errorf("%w: %s", err, out)
		return
	}

	return
}

func (k *Files) CopyFileFrom(localPath string, h RemoteFileHandle) (err error) {
	remotePath := path.Join(k.BasePath, h.FileName())

	cmd := exec.Command("cp", "-R", remotePath, localPath)
	out, err := cmd.CombinedOutput()

	if err != nil {
		err = xerrors.Errorf("%w: %s", err, out)
		return
	}

	return
}

// func (e *Files) GetAll() (zettels []string, err error) {
// 	glob := filepath.Join(e.BasePath, "*.md")
// 	zettels, err = filepath.Glob(glob)
// 	return
// }