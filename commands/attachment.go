package commands

import (
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"golang.org/x/xerrors"
)

type attachmentKind struct {
	hydrator func(u lib.Umwelt, urlString string) (z *lib.Zettel, err error)
}

func (a attachmentKind) String() string {
	//TODO-P4
	return ""
}

func (a *attachmentKind) Set(s string) (err error) {
	switch s {
	case "zettels-copy":
		*a = attachmentKind{
			hydrator: func(u lib.Umwelt, path string) (z *lib.Zettel, err error) {
				return pipeline.Import(u, path, true)
			},
		}
	case "zettels":
		*a = attachmentKind{
			hydrator: func(u lib.Umwelt, path string) (z *lib.Zettel, err error) {
				return pipeline.Import(u, path, false)
			},
		}
	case "files-copy":
		*a = attachmentKind{
			hydrator: func(u lib.Umwelt, urlString string) (z *lib.Zettel, err error) {
				return pipeline.NewOrFoundForFile(u, urlString, true)
			},
		}
	case "files":
		*a = attachmentKind{
			hydrator: func(u lib.Umwelt, urlString string) (z *lib.Zettel, err error) {
				return pipeline.NewOrFoundForFile(u, urlString, false)
			},
		}
	case "urls":
		*a = attachmentKind{
			hydrator: func(u lib.Umwelt, urlString string) (z *lib.Zettel, err error) {
				return pipeline.NewOrFoundForUrl(u, urlString)
			},
		}
	default:
		err = xerrors.Errorf("unsupported type: '%s'", s)
		return
	}

	return
}
