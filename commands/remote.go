package commands

import (
	"flag"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"golang.org/x/xerrors"
)

func init() {
	makeAndRegisterCommand(
		"remote",
		GetSubcommandRemote,
	)
}

func GetSubcommandRemote(f *flag.FlagSet) lib.Transactor {
	var query string
	f.StringVar(&query, "query", "", "zettel-spec")

	return func(u *lib.Umwelt) (err error) {
		args := f.Args()

		var command options.RemoteCommand

		err = command.Set(args[0])

		if err != nil {
			return
		}

		args = args[1:]

		if len(args) == 0 {
			err = xerrors.Errorf("remote name required")
			return
		}

		remoteName := args[0]

		remote, ok := u.Config.RemoteScripts[remoteName]

		if !ok {
			err = xerrors.Errorf("no remote with name '%s'", remoteName)
			return
		}

		remotePipeline, err := pipeline.MakeRemoteScript(u, &command, remote)

		if err != nil {
			return
		}

		args = args[1:]

		if len(args) != 0 {
			err = xerrors.Errorf("remote does not support positional arguments")
			return
		}

		remotePipeline.Run(u)

		return
	}
}
