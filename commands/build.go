package commands

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandBuild(f *flag.FlagSet) CommandRunFunc {
	return func(e *lib.Env) (err error) {
		actioner := func(i int, z *lib.Zettel) (shouldPrint bool, actionErr error) {
			shouldPrint = true

			var name string
			name, actionErr = getZettelBuildFileName(z)
			sPath := path.Join(e.BasePath, "build", name)

			if actionErr != nil {
				return
			}

			actionErr = syscall.Link(z.Path, sPath)

			if actionErr != nil && !os.IsExist(actionErr) {
				actionErr = fmt.Errorf("linking: %s: %w", sPath, actionErr)
				return
			}

			for _, t := range z.IndexData.Tags {
				actionErr = symlinkZettel(e, t, z)

				if actionErr != nil {
					actionErr = fmt.Errorf("symlinking zettel to tag: %w", actionErr)
					return
				}
			}

			return
		}

		processor := MakeProcessor(
			e,
			f.Args(),
			&nullZettelPrinter{},
		)

		processor.actioner = actioner

		buildPath := path.Join(e.BasePath, "build")

		os.RemoveAll(buildPath)
		err = os.Mkdir(buildPath, 0700)

		if err != nil && !os.IsExist(err) {
			return
		}

		err = processor.Run()

		return
	}
}

func symlinkZettel(e *lib.Env, dir string, z *lib.Zettel) (err error) {
	pPath, err := makeDirectoryIfNecessary(e, dir)

	if err != nil {
		err = fmt.Errorf("making directory: %s: %w", dir, err)
		return
	}

	pzPath, err := getZettelBuildFileName(z)

	if err != nil {
		err = fmt.Errorf("making zettel symlink: %s: %w", pPath, err)
		return
	}

	pzPath = path.Join(pPath, pzPath)
	err = syscall.Link(z.Path, pzPath)

	if os.IsExist(err) {
		err = nil
	} else if err != nil && !os.IsExist(err) {
		err = fmt.Errorf("linking zettel: %s: %w", pzPath, err)
		return
	}

	return
}

func makeDirectoryIfNecessary(e *lib.Env, p string) (a string, err error) {
	a = path.Join(e.BasePath, "build", p)
	err = os.Mkdir(a, 0700)

	if os.IsExist(err) {
		err = nil
	}

	return
}

func getZettelBuildFileName(z *lib.Zettel) (path string, err error) {
	sb := &strings.Builder{}
	t, err := lib.TimeFromPath(z.Path)

	if err != nil {
		return
	}

	day := t.Format("2006-01-02")

	sb.WriteString(day)
	sb.WriteString(" ")

	sb.WriteString(strings.ReplaceAll(z.IndexData.Description, "/", "-"))
	sb.WriteString(".md")

	path = sb.String()

	if len(path) > 255 {
		path = path[:255]
	}

	return
}
