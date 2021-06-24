package commands

import (
	"flag"
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandNew(f *flag.FlagSet) CommandRunFunc {
	return func(e Env) (err error) {
		currentTime := time.Now()
		unixTime := currentTime.Unix()

		filename := path.Join(e.ZettelPath, strconv.FormatInt(unixTime, 10)+".md")

		z := &lib.Zettel{
			Path: filename,
			Metadata: lib.ZettelMetadata{
				Date: currentTime.Format("2006-01-02"),
				Kind: "unknown",
				Tags: []string{"open"},
			},
		}

		err = z.Write()

		if err != nil {
			return
		}

		fmt.Print(filename)

		return
	}
}
