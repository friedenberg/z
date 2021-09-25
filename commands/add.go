package commands

import (
	"flag"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

type attachmentKind struct {
	adder    func(z *lib.KastenZettel, t time.Time, p string) lib.OnZettelWriteFunc
	hydrator func(u lib.Umwelt, urlString string) (z *lib.KastenZettel, err error)
}

func (a attachmentKind) String() string {
	//TODO
	return ""
}

func (a *attachmentKind) Set(s string) (err error) {
	fileAdder := func(z *lib.KastenZettel, t time.Time, p string) lib.OnZettelWriteFunc {
		z.Metadata.File = strconv.FormatInt(z.Id, 10) + path.Ext(p)
		return lib.AddFileOnWrite(p)
	}

	switch s {
	case "files-copy":
		*a = attachmentKind{
			adder: fileAdder,
			hydrator: func(u lib.Umwelt, urlString string) (z *lib.KastenZettel, err error) {
				return pipeline.NewOrFoundForFile(u, urlString, true)
			},
		}
	case "files":
		*a = attachmentKind{
			adder: fileAdder,
			hydrator: func(u lib.Umwelt, urlString string) (z *lib.KastenZettel, err error) {
				return pipeline.NewOrFoundForFile(u, urlString, false)
			},
		}
	case "urls":
		*a = attachmentKind{
			adder: func(z *lib.KastenZettel, t time.Time, p string) lib.OnZettelWriteFunc {
				//TODO normalize
				z.Metadata.Url = p
				return lib.AddUrlOnWrite(p, t)
			},
			hydrator: func(u lib.Umwelt, urlString string) (z *lib.KastenZettel, err error) {
				return pipeline.NewOrFoundForUrl(u, urlString)
			},
		}
	default:
		err = xerrors.Errorf("unsupported type: '%s'", s)
		return
	}

	return
}

func GetSubcommandAdd(f *flag.FlagSet) CommandRunFunc {
	var tagString string
	var kind attachmentKind
	editActions := options.Actions(options.ActionEdit)

	f.Var(&editActions, "actions", "action to perform for the matched zettels")

	f.StringVar(&tagString, "tags", "", "parse the passed-in string as the metadata.")
	f.Var(&kind, "kind", "treat the positional arguments as this kind.")

	return func(e lib.Umwelt) (err error) {
		currentTime := time.Now()

		pr := &printer.MultiplexingZettelPrinter{
			Printer: &printer.ActionZettelPrinter{
				Umwelt:  e,
				Actions: editActions,
			},
		}

		iter := func(i int, a string) (err error) {
			z, err := kind.hydrator(e, a)

			if err != nil {
				return
			}

			if tagString != "" {
				tags := strings.Split(tagString, " ")
				z.Metadata.Tags = append(z.Metadata.Tags, tags...)
			} else {
				z.Metadata.Tags = append(z.Metadata.Tags, "zz-inbox")
			}

			z.Metadata.Tags = uniqueAndSortTags(z.Metadata.Tags)

			onWrite := kind.adder(z, currentTime, a)

			err = z.Write(onWrite)

			if err != nil {
				err = xerrors.Errorf("failed to write: %w", err)
			}

			return
		}

		par := util.Parallelizer{Args: f.Args()}
		pr.Printer.Begin()
		defer pr.Printer.End()
		par.Run(iter, errIterartion(pr.Printer))

		return
	}
}
