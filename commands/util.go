package commands

import (
	"sort"

	"github.com/friedenberg/z/commands/printer"
	"github.com/friedenberg/z/lib"
)

func hydrateIndex(k *lib.FilesAndGit) error {
	indexProcessor := MakeProcessor(
		k,
		nil,
		&printer.NullZettelPrinter{},
	)

	indexProcessor.actioner = func(i int, z *lib.Zettel) (shouldPrint bool, err error) {
		err = k.Index.Add(z)
		return
	}

	return indexProcessor.Run()
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
