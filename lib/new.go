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
)

func (z *Zettel) InitFromTime(t time.Time) {
	z.Path = MakePathFromTime(z.Env.BasePath, t)
	z.Id = t.Unix()

	z.IndexData = ZettelIndexData{
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

		z.IndexData.Url = url.String()

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

		z.Data.Body = md.String()

		return
	}
}

func AddFileOnWrite(oldPath string) OnZettelWriteFunc {
	return func(z *Zettel, errIn error) (errOut error) {

		if errIn != nil {
			return
		}

		errOut = os.Rename(oldPath, z.FilePath())

		if errOut != nil {
			errOut = fmt.Errorf("rename file: %w", errOut)
		}

		return
	}
}
