package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func GetSubcommandCat(f *flag.FlagSet) CommandRunFunc {
	var shouldPrintJson bool
	f.BoolVar(&shouldPrintJson, "json", false, "TODO")

	return func(e Env) (err error) {
		var putter Putter

		if shouldPrintJson {
			putter = &jsonPutter{channel: make(PutterChannel)}
		} else {
			putter = &basicPutter{channel: make(PutterChannel)}
		}

		files := f.Args()

		if len(files) == 0 {
			glob := filepath.Join(e.ZettelPath, "*.md")
			files, err = filepath.Glob(glob)

			if err != nil {
				return
			}
		}

		processor := MakeProcessor(
			e,
			files,
			putter,
		)

		err = processor.Run()

		return
	}
}

type basicPutter struct {
	channel PutterChannel
}

func (p *basicPutter) GetChannel() PutterChannel {
	return p.channel
}

func (p *basicPutter) Print() {
	for z := range p.channel {
		fmt.Println(z.Data.MetadataYaml)
		fmt.Println(z.Data.Body)
	}
}

type jsonPutter struct {
	channel PutterChannel
}

func (p *jsonPutter) GetChannel() PutterChannel {
	return p.channel
}

func (p *jsonPutter) Print() {
	for z := range p.channel {
		j, err := json.Marshal(z.IndexData)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		fmt.Println(string(j))
	}
}
