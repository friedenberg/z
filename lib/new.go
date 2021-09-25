package lib

import (
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

//TODO deprecate
func (z *KastenZettel) InitAndAssignUniqueId(currentTime time.Time, i int) (err error) {
	d, err := time.ParseDuration(strconv.Itoa(i) + "s")

	if err != nil {
		return
	}

	t := currentTime.Add(d)
	z.InitFromTime(t)

	for {
		if util.FileExists(z.Path) {
			d, err = time.ParseDuration("1s")

			if err != nil {
				return
			}

			currentTime = currentTime.Add(d)
			z.InitFromTime(currentTime)
		} else {
			break
		}
	}

	return
}

func (z *KastenZettel) InitFromTime(t time.Time) {
	z.Path = MakePathFromTime(z.Umwelt.FilesAndGit().BasePath, t)
	z.Id = t.Unix()

	z.Metadata = Metadata{
		Date: t.Format("2006-01-02"),
	}

	return
}

func MakePathFromId(basePath, id string) string {
	return path.Join(basePath, id+".md")
}

func MakePathFromTime(basePath string, t time.Time) (filename string) {
	unixTime := t.Unix()
	return MakePathFromId(basePath, strconv.FormatInt(unixTime, 10))
}

func AddUrlOnWrite(u string, t time.Time) OnZettelWriteFunc {
	return func(z *KastenZettel, errIn error) (errOut error) {
		if errIn != nil {
			return
		}

		url, errOut := url.Parse(u)

		if errOut != nil {
			return
		}

		z.Metadata.Url = url.String()

		//TODO determine summaries from sites
		return
	}
}

func AddFileOnWrite(oldPath string) OnZettelWriteFunc {
	return func(z *KastenZettel, errIn error) (errOut error) {

		if errIn != nil {
			return
		}

		errOut = os.Rename(oldPath, z.FilePath())
		// cmd := exec.Command("cp", "-R", oldPath, z.FilePath())
		// errOut = cmd.Run()

		if errOut != nil {
			errOut = xerrors.Errorf("cp file: %w", errOut)
		}

		return
	}
}
