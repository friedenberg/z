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
	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

type subcommand struct {
	flags   *flag.FlagSet
	runFunc commands.CommandRunFunc
}

var (
	subcommands = map[string]subcommand{}
)

func makeSubcommand(name string, makeFunc func(*flag.FlagSet) commands.CommandRunFunc) {
	if _, ok := subcommands[name]; ok {
		panic("command added more than once: " + name)
	}

	flags := flag.NewFlagSet(name, flag.ExitOnError)
	subcommands[name] = subcommand{
		flags:   flags,
		runFunc: makeFunc(flags),
	}
}

func init() {
	makeSubcommand("add", commands.GetSubcommandAdd)
	makeSubcommand("autocomplete", commands.GetSubcommandAutocomplete)
	makeSubcommand("build", commands.GetSubcommandBuild)
	makeSubcommand("cat", commands.GetSubcommandCat)
	makeSubcommand("clean", commands.GetSubcommandClean)
	makeSubcommand("edit", commands.GetSubcommandEdit)
	makeSubcommand("index", commands.GetSubcommandIndex)
	makeSubcommand("mv", commands.GetSubcommandMv)
	makeSubcommand("new", commands.GetSubcommandNew)
	makeSubcommand("remote", commands.GetSubcommandRemote)
	makeSubcommand("rm", commands.GetSubcommandRm)
}

func main() {
	os.Exit(run())
}

func run() int {
	defer util.WaitForPrinter()

	if len(os.Args) < 2 {
		return printUsage(nil)
	}

	specifiedSubcommand := os.Args[1]
	cmd, ok := subcommands[specifiedSubcommand]

	if !ok {
		return printUsage(xerrors.Errorf("No subcommand '%s'", specifiedSubcommand))
	}

	var err error
	c, err := lib.LoadDefaultConfig()

	if err != nil {
		util.StdPrinterErr(err)
		return 1
	}

	umwelt, err := c.Umwelt()

	if err != nil {
		util.StdPrinterError(err)
		return 1
	}

	//TODO refactor to be command too
	cmd.flags.Parse(os.Args[2:])
	err = cmd.runFunc(umwelt)

	if err != nil {
		util.StdPrinterError(err)
	}

	return 0
}

func printUsage(err error) int {
	if err != nil {
		util.StdPrinterErr(err)
	}

	fmt.Println("Usage for z: ")

	sc := make([]subcommand, 0, len(subcommands))

	for _, c := range subcommands {
		sc = append(sc, c)
	}

	sort.Slice(sc, func(i, j int) bool {
		return sc[i].flags.Name() < sc[j].flags.Name()
	})

	for _, c := range sc {
		printSubcommandUsage(c)
	}

	status := 0

	if err != nil {
		//TODO get correct status
		status = 1
	}

	return status
}

func printSubcommandUsage(sc subcommand) {
	printTabbed := func(s string) {
		util.StdPrinterErrf("  %s\n", s)
	}

	flags := sc.flags

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
