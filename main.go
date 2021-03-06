package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/friedenberg/z/commands"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/util/stdprinter"
	"golang.org/x/xerrors"
)

func main() {
	os.Exit(run())
}

func run() int {
	stdprinter.SetDebug(shouldDebug())
	defer stdprinter.WaitForPrinter()

	if len(os.Args) < 1 {
		stdprinter.Debug("printing usage")
		return printUsage(nil)
	}

	cmds := commands.Commands()
	specifiedSubcommand := os.Args[1]
	cmd, ok := cmds[specifiedSubcommand]

	if !ok {
		return printUsage(xerrors.Errorf("No subcommand '%s'", specifiedSubcommand))
	}

	var err error
	c, err := lib.LoadDefaultConfig()

	if err != nil {
		stdprinter.Err(err)
		return 1
	}

	umwelt, err := c.Umwelt()

	if err != nil {
		stdprinter.Err(err)
		return 1
	}

	//TODO-P4 refactor to be command too
	cmd.Flags.Parse(os.Args[2:])
	err = cmd.Run(umwelt)

	if err != nil {
		stdprinter.Error(err)
		return 1
	}

	umwelt.Transaction.IsFinalTransaction = true

	err = umwelt.Kasten.CommitTransaction(umwelt)

	if err != nil {
		stdprinter.Error(err)
		return 1
	}

	return 0
}

func printUsage(err error) int {
	if err != nil {
		stdprinter.Err(err)
	}

	fmt.Println("Usage for z: ")

	fs := make([]flag.FlagSet, 0, len(commands.Commands()))

	for _, c := range commands.Commands() {
		fs = append(fs, *c.Flags)
	}

	sort.Slice(fs, func(i, j int) bool {
		return fs[i].Name() < fs[j].Name()
	})

	for _, f := range fs {
		printSubcommandUsage(f)
	}

	status := 0

	if err != nil {
		//TODO-P5 get correct status
		status = 1
	}

	return status
}

func printSubcommandUsage(flags flag.FlagSet) {
	printTabbed := func(s string) {
		stdprinter.Errf("  %s\n", s)
	}

	var b bytes.Buffer
	flags.SetOutput(&b)

	printTabbed(flags.Name())
	flags.PrintDefaults()

	scanner := bufio.NewScanner(&b)

	for scanner.Scan() {
		printTabbed(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func shouldDebug() bool {
	for i, a := range os.Args {
		if a == "-v" {
			stdprinter.Debug("debugging on")
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			return true
		}
	}

	return false
}
