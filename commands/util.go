package commands

import (
	"sort"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/pipeline/printer"
	"github.com/friedenberg/z/util"
)

func hydrateIndex(k lib.Umwelt) (err error) {
	allFiles, err := k.Store().GetAll()

	if err != nil {
		return
	}

	indexProcessor := MakeProcessor(
		k,
		allFiles,
		&printer.NullZettelPrinter{},
	)

	indexProcessor.hydrator = HydrateFromFileFunc(k, true)

	indexProcessor.actioner = func(i int, z *lib.Zettel) (shouldPrint bool, err error) {
		err = k.Index.Add(z)
		return
	}

	err = indexProcessor.Run()

	if err != nil {
		return
	}

	return
}

func uniqueAndSortTags(tags []string) (o []string) {
	o = make([]string, 0, len(tags))
	m := make(map[string]bool)

	for _, t := range tags {
		if _, ok := m[t]; !ok {
			m[t] = true
			o = append(o, t)
		}
	}

	sort.Slice(o, func(i, j int) bool {
		return o[i] < o[j]
	})

	return
}

//TODO-P4 refactor
func doesZettelMatchQuery(z *lib.Zettel, q string) bool {
	if q == "" {
		return true
	}

	//TODO-P2
	// if z.Note.Metadata.LocalFile() == q {
	// 	return true
	// }

	// if z.Metadata.Url == q {
	// 	return true
	// }

	for _, t := range z.Note.Metadata.TagStrings() {
		if t == q {
			return true
		}
	}

	return false
}

func errIterartion(p pipeline.Printer) util.ParallelizerErrorFunc {
	return func(i int, s string, err error) {
		p.PrintZettel(i, nil, err)
	}
}

func printIfNecessary(i int, z *lib.Zettel, q string, fp pipeline.FilterPrinter) {
	if (fp.Filter == nil || fp.Filter(i, z)) && doesZettelMatchQuery(z, q) {
		fp.Printer.PrintZettel(i, z, nil)
	}
}

func cachedIteration(u lib.Umwelt, q string, fp pipeline.FilterPrinter) util.ParallelizerIterFunc {
	return func(i int, s string) (err error) {
		s = util.BaseNameNoSuffix(s)
		z, err := pipeline.HydrateFromIndex(u, s)

		if err != nil {
			return
		}

		printIfNecessary(i, z, q, fp)

		return
	}
}

func filesystemIteration(u lib.Umwelt, q string, fp pipeline.FilterPrinter) util.ParallelizerIterFunc {
	return func(i int, s string) (err error) {
		p, err := pipeline.NormalizePath(u, s)
		p = util.EverythingExceptExtension(p) + ".md"

		if err != nil {
			return
		}

		z, err := pipeline.HydrateFromFile(u, p, true)

		if err != nil {
			return
		}

		printIfNecessary(i, z, q, fp)

		return
	}
}
