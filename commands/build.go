package commands

import (
	"flag"
	"os"
	"path"
	"strings"
	"syscall"

	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
	"golang.org/x/xerrors"
)

func GetSubcommandBuild(f *flag.FlagSet) CommandRunFunc {
	return func(e lib.Umwelt) (err error) {
		actioner := func(i int, z *lib.KastenZettel) (shouldPrint bool, actionErr error) {
			shouldPrint = true

			for _, t := range z.Metadata.Tags {
				actionErr = symlinkZettel(e, t, z)

				if actionErr != nil {
					actionErr = xerrors.Errorf("symlinking zettel to tag: %w", actionErr)
					return
				}
			}

			if len(z.Metadata.Tags) == 0 {
				actionErr = symlinkZettel(e, "untagged", z)

				if actionErr != nil {
					actionErr = xerrors.Errorf("symlinking zettel: %w", actionErr)
					return
				}
			}

			return
		}

		processor := MakeProcessor(
			e,
			f.Args(),
			&printer.NullZettelPrinter{},
		)

		processor.actioner = actioner

		buildPath := path.Join(e.FilesAndGit().BasePath, "build")

		os.RemoveAll(buildPath)
		err = os.Mkdir(buildPath, 0700)

		if err != nil && !os.IsExist(err) {
			return
		}

		err = processor.Run()

		return
	}
}

func symlinkZettel(e lib.Umwelt, dir string, z *lib.KastenZettel) (err error) {
	buildDir, err := makeDirectoryIfNecessary(e, dir)

	if err != nil {
		err = xerrors.Errorf("making directory: %s: %w", dir, err)
		return
	}

	doSym := func(originalPath, ext string) error {
		newFilename, err := getZettelBuildFileName(z, ext)
		// fmt.Println(originalPath)
		// fmt.Println(newFilename)

		if err != nil {
			return xerrors.Errorf("making zettel symlink: %s: %w", originalPath, err)
		}

		symPath := path.Join(buildDir, newFilename)
		err = syscall.Link(originalPath, symPath)

		if err != nil && !os.IsExist(err) {
			return xerrors.Errorf("linking zettel: %s: %w", symPath, err)
		}

		return nil
	}

	err = doSym(z.Path, ".md")

	if err != nil {
		return
	}

	if z.HasFile() {
		err = doSym(z.FilePath(), path.Ext(z.Metadata.File))
	}

	return
}

func makeDirectoryIfNecessary(e lib.Umwelt, p string) (a string, err error) {
	a = path.Join(e.FilesAndGit().BasePath, "build", p)
	err = os.Mkdir(a, 0700)

	if os.IsExist(err) {
		err = nil
	}

	return
}

func getZettelBuildFileName(z *lib.KastenZettel, ext string) (path string, err error) {
	sb := &strings.Builder{}
	t, err := lib.TimeFromPath(z.Path)

	if err != nil {
		return
	}

	day := t.Format("2006-01-02")

	sb.WriteString(day)
	sb.WriteString(" ")

	sb.WriteString(strings.ReplaceAll(z.Metadata.Description, "/", "-"))
	sb.WriteString(ext)

	path = sb.String()

	if len(path) > 255 {
		path = path[:255]
	}

	return
}
