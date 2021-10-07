package commands

import (
	"flag"
	"strings"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/pipeline/printer"
	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

type attachmentKind struct {
	hydrator func(u lib.Umwelt, urlString string) (z *lib.Zettel, err error)
}

func (a attachmentKind) String() string {
	//TODO
	return ""
}

func (a *attachmentKind) Set(s string) (err error) {
	switch s {
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

func GetSubcommandAdd(f *flag.FlagSet) CommandRunFunc {
	var tagString string
	var description string
	var kind attachmentKind
	editActions := options.Actions(options.ActionEdit)

	f.Var(&editActions, "actions", "action to perform for the matched zettels")

	f.StringVar(&tagString, "tags", "", "parse the passed-in string as the metadata.")
	f.StringVar(&description, "description", "", "use this string as the zettel description")
	f.Var(&kind, "kind", "treat the positional arguments as this kind.")

	return func(e lib.Umwelt) (err error) {
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

			z.Metadata.Description = description

			z.Metadata.Tags = uniqueAndSortTags(z.Metadata.Tags)

			err = z.Write(nil)

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
