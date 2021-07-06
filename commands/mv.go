package commands

import (
	"flag"
	"fmt"
	"strings"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandMv(f *flag.FlagSet) CommandRunFunc {
	isDryRun := false

	f.BoolVar(&isDryRun, "dry-run", false, "")

	return func(e *lib.Env) (err error) {
		args := f.Args()

		fromMoveInstruction, err := moveInstructionFromString(args[0])
		if err != nil {
			return
		}

		toMoveInstruction, err := moveInstructionFromString(args[1])

		if err != nil {
			return
		}

		processor := MakeProcessor(
			e,
			args[2:],
			&nullZettelPrinter{},
		)

		processor.actioner = func(_ int, z *lib.Zettel) (err error) {
			found := -1
			values := fromMoveInstruction.fieldReadWriter.ValueGetFunc(z)

			for i, v := range values {
				if v == fromMoveInstruction.value {
					found = i
					break
				}
			}

			if found < 0 {
				return
			}

			values = append(values[:found], values[found+1:]...)
			fromMoveInstruction.fieldReadWriter.ValueSetFunc(z, values)

			values = toMoveInstruction.fieldReadWriter.ValueGetFunc(z)
			values = append(values, toMoveInstruction.value)
			toMoveInstruction.fieldReadWriter.ValueSetFunc(z, values)

			fmt.Println(z.IndexData)

			if !isDryRun {
				err = z.Write(func(_ *lib.Zettel, _ error) error { return nil })
			}

			if err != nil {
				err = fmt.Errorf("failed to write: %w", err)
				return
			}

			return
		}

		err = processor.Run()

		return
	}
}

type moveInstruction struct {
	fieldReadWriter lib.MetadataFieldReadWriterArray
	value           string
}

func moveInstructionFromString(s string) (m moveInstruction, err error) {
	if s == "" {
		m = moveInstruction{
			fieldReadWriter: lib.GetMetadataFieldReadWriterNull(),
		}

		return
	}

	split := strings.Split(s, ":")

	if len(split) != 2 {
		err = fmt.Errorf("'%s': incorrect number of field specifier characters (':')", s)
		return
	}

	fieldShortName := split[0]

	m = moveInstruction{
		value: split[1],
	}

	switch fieldShortName {
	case "a":
		m.fieldReadWriter = lib.GetMetadataFieldReadWriterAreas()
	case "p":
		m.fieldReadWriter = lib.GetMetadataFieldReadWriterProjects()
	case "t":
		m.fieldReadWriter = lib.GetMetadataFieldReadWriterTags()
	default:
		err = fmt.Errorf("'%s': invalid field short name", fieldShortName)
		return
	}

	return
}
