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

func makeCommand(n string, makeFunc func(*flag.FlagSet) lib.Transactor) (c Command) {
	f := flag.NewFlagSet(n, flag.ExitOnError)

	c = Command{
		Flags: f,
		Run:   makeFunc(f),
	}

	return
}

func makeAndRegisterCommand(n string, makeFunc func(*flag.FlagSet) lib.Transactor) {
	c := makeCommand(n, makeFunc)
	registerCommand(n, c)

	return
}

func Commands() map[string]Command {
	return commands
}
