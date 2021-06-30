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
	"time"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandAdd(f *flag.FlagSet) CommandRunFunc {
	isUrl := false
	f.BoolVar(&isUrl, "url", false, "")

	return func(e Env) (err error) {
		currentTime := time.Now()

		processor := MakeProcessor(
			e,
			f.Args(),
			&NullPutter{Channel: make(PutterChannel)},
		)

		processor.hydrateAction = func(i int, z *lib.Zettel) (err error) {
			d, err := time.ParseDuration(strconv.Itoa(i) + "s")
			t := currentTime.Add(d)
			if err != nil {
				panic(err)
			}
			unixTime := t.Unix()
			zettelId := strconv.FormatInt(unixTime, 10)

			zFilename := path.Join(e.ZettelPath, zettelId+".md")
			fmt.Println(zFilename)

			p := z.Path

			z.Path = zFilename
			z.Metadata = lib.ZettelMetadata{
				Date: t.Format("2006-01-02"),
				Tags: []string{"added"},
			}

			var onWrite onZettelWrite

			if isUrl {
				onWrite, err = addUrl(z, e, p, t)
			} else {
				onWrite, err = addFile(z, e, p, zettelId)
			}

			if err != nil {
				err = fmt.Errorf("failed to add url or file: %w", err)
				return
			}

			err = z.Write()

			if onWrite != nil {
				err = onWrite(z, err)
			}

			if err != nil {
				err = fmt.Errorf("failed to write: %w", err)
			}

			return
		}

		err = processor.Run()

		return
	}
}

type onZettelWrite func(*lib.Zettel, error) error

func addUrl(z *lib.Zettel, e Env, u string, t time.Time) (onWrite onZettelWrite, err error) {
	url, err := url.Parse(u)
	fmt.Println(url)
	fmt.Println(err)
	os.Exit(1)

	if err != nil {
		return
	}

	z.Metadata.Kind = "pb"
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

func addFile(z *lib.Zettel, e Env, p string, zi string) (onWrite onZettelWrite, err error) {
	newFilename := path.Join(e.ZettelPath, zi+path.Ext(p))
	z.Metadata.File = newFilename
	z.Metadata.Kind = "file"
	onWrite = func(z *lib.Zettel, err error) error {
		return os.Rename(p, newFilename)
	}
	return
}
