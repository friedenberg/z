package pipeline

import (
	"bufio"
	"flag"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/friedenberg/z/lib"
	"github.com/friedenberg/z/lib/feeder"
	"github.com/friedenberg/z/lib/zettel/filter"
	"github.com/friedenberg/z/util/stdprinter"
	"golang.org/x/xerrors"
)

type RemoteScript interface {
	Run(*lib.Umwelt) error
}

type remoteScript struct {
	Pipeline
	lib.RemoteScript
}

type RemoteScriptPull remoteScript
type RemoteScriptPush remoteScript

func MakeRemoteScript(u *lib.Umwelt, c1 flag.Value, c2 lib.ConfigRemoteScript) (s RemoteScript, err error) {
	var rs lib.RemoteScript

	switch c1.String() {
	case "pull":
		rs = c2.Pull
	case "push":
		rs = c2.Push
	default:
		err = xerrors.Errorf("unsupported remote script command: '%s'", c1.String())
		return
	}

	if rs.Script == "" {
		err = xerrors.Errorf("script body is empty")
		return
	}

	if rs.Shell == "" {
		rs.Shell = "bash"
	}

	f := Format{}
	err = f.Set(rs.Format)

	if err != nil {
		return
	}

	// start with using the specific filter for the script
	var f1 filter.Filter
	f1 = filter.Tag(rs.Filter)

	// then, include the filter for the format if it exists
	if f.Filter != nil {
		f1 = filter.MakeAnd(f1, f.Filter)
	}

	switch c1.String() {
	case "pull":
		s = &RemoteScriptPull{
			Pipeline: Pipeline{
				Reader:   f.Reader,
				Filter:   f1,
				Modifier: lib.MakeTransactionAction(u.Transaction, lib.TransactionActionAdded),
			},
			RemoteScript: rs,
		}
	case "push":
		s = &RemoteScriptPush{
			Pipeline: Pipeline{
				Filter: f1,
				//TODO add modifier to committing remote push
				Writer: f.Writer,
			},
			RemoteScript: rs,
		}
	default:
		panic("not possible")
		return
	}

	return
}

func (s *RemoteScriptPull) Run(u *lib.Umwelt) (err error) {
	script := exec.Command(s.Shell, "-c", s.Script)

	script.Stderr = os.Stderr

	r, w := io.Pipe()
	script.Stdout = w

	script.Start()

	args := make([]string, 0, 0)
	sr := bufio.NewReader(r)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		script.Wait()
		w.Close()
		stdprinter.Debug("done waiting")
	}()

	go func() {
		defer wg.Done()
		for {
			var s string
			s, err = sr.ReadString('\n')

			if err == io.EOF {
				stdprinter.Debug("eof", s)
				err = nil
				break
			}

			if err != nil {
				return
			}

			args = append(args, s)
			stdprinter.Debug("read:", s)
		}
	}()

	wg.Wait()

	s.Pipeline.Feeder = feeder.MakeStringSlice(args)
	s.Pipeline.Out = os.Stdout
	s.Pipeline.Run(u)

	return
}

func (s *RemoteScriptPush) Run(u *lib.Umwelt) (err error) {
	script := exec.Command(s.Shell, "-c", s.Script)

	script.Stdout = os.Stdout
	script.Stderr = os.Stderr

	r, w := io.Pipe()
	script.Stdin = r

	script.Start()

	s.Pipeline.Feeder = u.GetAll()
	s.Pipeline.Out = w
	s.Pipeline.Run(u)

	w.Close()
	script.Wait()

	return
}
