package lib

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/friedenberg/z/util"
)

func (z *Zettel) InitAndAssignUniqueId(currentTime time.Time, i int) (err error) {
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

func (z *Zettel) InitFromTime(t time.Time) {
	z.Path = MakePathFromTime(z.Kasten.BasePath, t)
	z.Id = t.Unix()

	z.Metadata = Metadata{
		Date: t.Format("2006-01-02"),
	}

	return
}

func MakePathFromTime(basePath string, t time.Time) (filename string) {
	unixTime := t.Unix()
	filename = path.Join(basePath, strconv.FormatInt(unixTime, 10)+".md")

	return
}

func AddUrlOnWrite(u string, t time.Time) OnZettelWriteFunc {
	return func(z *Zettel, errIn error) (errOut error) {
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

		//TODO description from title

		chromeCommand := exec.Command(
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"--headless",
			"--dump-dom",
			url.String(),
		)

		//TODO add filters
		pandocCommand := exec.Command(
			"pandoc",
			"-f",
			"html-native_divs-native_spans",
			"-t",
			"markdown-header_attributes-fenced_code_attributes-link_attributes",
		)

		r, w := io.Pipe()
		chromeCommand.Stdout = w
		pandocCommand.Stdin = r

		var md bytes.Buffer
		pandocCommand.Stdout = &md

		chromeCommand.Start()
		pandocCommand.Start()
		chromeCommand.Wait()
		w.Close()
		pandocCommand.Wait()

		z.Body = md.String()

		return
	}
}

func AddFileOnWrite(oldPath string) OnZettelWriteFunc {
	return func(z *Zettel, errIn error) (errOut error) {

		if errIn != nil {
			return
		}

		errOut = os.Rename(oldPath, z.FilePath())
		// cmd := exec.Command("cp", "-R", oldPath, z.FilePath())
		// errOut = cmd.Run()

		if errOut != nil {
			errOut = fmt.Errorf("cp file: %w", errOut)
		}

		return
	}
}
