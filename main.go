package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/friedenberg/z/commands"
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
	makeSubcommand("cat", commands.GetSubcommandCat)
	makeSubcommand("clean", commands.GetSubcommandClean)
	makeSubcommand("edit", commands.GetSubcommandEdit)
	makeSubcommand("new", commands.GetSubcommandNew)
	makeSubcommand("open", commands.GetSubcommandOpen)
	makeSubcommand("print", commands.GetSubcommandPrint)
	makeSubcommand("rm", commands.GetSubcommandRm)
}

func main() {
	var err error
	defaultEnv, err := commands.GetDefaultEnv()

	if err != nil {
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		printUsage(nil)
	}

	specifiedSubcommand := os.Args[1]
	cmd, ok := subcommands[specifiedSubcommand]

	if !ok {
		printUsage(fmt.Errorf("No subcommand '%s'", specifiedSubcommand))
	}

	cmd.flags.Parse(os.Args[2:])
	err = cmd.runFunc(defaultEnv)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func printUsage(err error) {
	if err != nil {
		fmt.Println(err)
		fmt.Println()
	}

	fmt.Println("Usage for z: ")

	for _, c := range subcommands {
		printSubcommandUsage(c)
	}

	status := 0

	if err != nil {
		//TODO get correct status
		status = 1
	}

	os.Exit(status)
}

func printSubcommandUsage(sc subcommand) {
	printTabbed := func(s string) {
		fmt.Println("  ", s)
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
