package commands

import (
	"flag"
	"io"
	"os"
	"os/exec"

	"github.com/friedenberg/z/commands/options"
	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/pipeline"
	"github.com/friedenberg/z/lib/zettel/filter"
	"github.com/friedenberg/z/lib/zettel/writer"
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

	return func(u lib.Umwelt) (err error) {
		args := f.Args()

		var command options.RemoteCommand

		err = command.Set(args[0])

		if err != nil {
			return
		}

		args = args[1:]

		remoteName := args[0]

		remote, ok := u.Config.RemoteScripts[remoteName]

		if !ok {
			err = xerrors.Errorf("no remote with name '%s'", remoteName)
			return
		}

		remotePipeline := remote.RemoteScriptForCommand(&command)

		if remotePipeline == nil {
			err = xerrors.Errorf("remote '%s' does not support command '%s'", remoteName, command)
			return
		}

		var format writer.Format
		err = format.Set(remotePipeline.Format)

		if err != nil {
			return
		}

		//TODO-P3 validation
		args = args[1:]

		if len(args) == 0 {
			args = u.GetAll()
		}

		script := exec.Command(
			remotePipeline.Shell,
			"-c",
			remotePipeline.Script,
		)

		script.Stdout = os.Stdout
		script.Stderr = os.Stderr

		r, w := io.Pipe()
		script.Stdin = r

		script.Start()

		p := pipeline.Pipeline{
			Arguments: args,
			Filter: filter.And(
				filter.MatchQuery(remotePipeline.Filter),
				format.Filter,
			),
			Writer: format.Writer,
			Out:    w,
		}

		p.Run(u)

		w.Close()
		script.Wait()

		return
	}
}
