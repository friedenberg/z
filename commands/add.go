package commands

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandAdd(f *flag.FlagSet) CommandRunFunc {
	return func(e Env) (err error) {
		currentTime := time.Now()
		wg := sync.WaitGroup{}

		for i, p := range f.Args() {
			d, err := time.ParseDuration(strconv.Itoa(i) + "s")

			if err != nil {
				panic(err)
			}

			t := currentTime.Add(d)
			wg.Add(1)
			go func() {
				defer wg.Done()
				unixTime := t.Unix()
				zettelId := strconv.FormatInt(unixTime, 10)

				zFilename := path.Join(e.ZettelPath, zettelId+".md")
				fmt.Println(zFilename)

				z := &lib.Zettel{
					Path: zFilename,
					Metadata: lib.ZettelMetadata{
						Date: t.Format("2006-01-02"),
						Tags: []string{"open"},
					},
				}

				onWrite, err := addUrl(z, e, p, t)

				if err != nil {
					onWrite, err = addFile(z, e, p, zettelId)
				}

				err = z.Write()

				if onWrite != nil {
					err = onWrite(z, err)
				}
				//TODO
			}()
		}

		wg.Wait()

		return
	}
}

type onZettelWrite func(*lib.Zettel, error) error

func addUrl(z *lib.Zettel, e Env, u string, t time.Time) (onWrite onZettelWrite, err error) {
	url, err := url.Parse(u)

	z.Metadata.Kind = "pb"
	//TODO description from title

	if err != nil {
		return
	}

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

func addFile(z *lib.Zettel, e Env, p string, zi string) (onWrite onZettelWrite, err error) {
	newFilename := path.Join(e.ZettelPath, zi+path.Ext(p))
	z.Metadata.File = p
	z.Metadata.Kind = "file"
	onWrite = func(z *lib.Zettel, err error) error {
		return os.Rename(p, newFilename)
	}
	return
}
