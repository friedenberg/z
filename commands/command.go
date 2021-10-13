package commands

import (
	"flag"

	"github.com/friedenberg/z/lib"
)

type Command struct {
	Flags *flag.FlagSet
	Run   lib.Transactor
}

var (
	commands = map[string]Command{}
)

func registerCommand(n string, c Command) {
	if _, ok := commands[n]; ok {
		panic("command added more than once: " + n)
	}

	commands[n] = c
}

func Commands() map[string]Command {
	return commands
}

func registerSubcommand(name string, makeFunc func(*flag.FlagSet) lib.Transactor) {
	flags := flag.NewFlagSet(name, flag.ExitOnError)
	s := Command{
		Flags: flags,
		Run:   makeFunc(flags),
	}

	registerCommand(name, s)
}

func init() {
	registerSubcommand("autocomplete", GetSubcommandAutocomplete)
	registerSubcommand("build", GetSubcommandBuild)
	registerSubcommand("cat", GetSubcommandCat)
	registerSubcommand("clean", GetSubcommandClean)
	registerSubcommand("edit", GetSubcommandEdit)
	registerSubcommand("index", GetSubcommandIndex)
	registerSubcommand("mv", GetSubcommandMv)
	registerSubcommand("new", GetSubcommandNew)
	registerSubcommand("remote", GetSubcommandRemote)
	registerSubcommand("rm", GetSubcommandRm)
}
