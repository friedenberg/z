package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandRm(f *flag.FlagSet) CommandRunFunc {
	return func(e Env) (err error) {
		path := f.Arg(0)

		if path == "" {
			err = errors.New("path was empty")
		}

		absPath, err := filepath.Abs(path)

		if err != nil {
			err = fmt.Errorf("%s: get absolute path: %w", path, err)
			return
		}

		z := &lib.Zettel{Path: absPath}
		z.HydrateFromFilePath()

		err = os.Remove(z.Path)

		if err != nil {
			return
		}

		if z.Metadata.Kind == "file" {
			err = os.Remove(z.Metadata.File)
		}

		return
	}
}
