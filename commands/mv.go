package commands

import (
	"flag"
	"fmt"

	"github.com/friedenberg/z/commands/printer"
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
			&printer.NullZettelPrinter{},
		)

		processor.hydrator = func(_ int, z *lib.Zettel, path string) (err error) {
			z.Path = path
			return z.HydrateFromFilePath(true)
		}

		processor.actioner = func(_ int, z *lib.Zettel) (shouldPrint bool, err error) {
			shouldPrint = true
			found := -1
			values := fromMoveInstruction.fieldReadWriter.ValueGetFunc(z)

			for i, v := range values {
				if v == fromMoveInstruction.value {
					found = i
					break
				}
			}

			if found < 0 {
				shouldPrint = false
				return
			}

			values = append(values[:found], values[found+1:]...)
			fromMoveInstruction.fieldReadWriter.ValueSetFunc(z, values)

			values = toMoveInstruction.fieldReadWriter.ValueGetFunc(z)
			values = append(values, toMoveInstruction.value)
			toMoveInstruction.fieldReadWriter.ValueSetFunc(z, values)

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

	m = moveInstruction{
		value:           s,
		fieldReadWriter: lib.GetMetadataFieldReadWriterTags(),
	}

	return
}
