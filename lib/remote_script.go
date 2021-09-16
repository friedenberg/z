package lib

import (
	"flag"
)

type RemoteScript struct {
	Format string
	Filter string
	Shell  string
	Script string
}

func (c ConfigRemoteScript) RemoteScriptForCommand(co flag.Value) (s *RemoteScript) {
	switch co.String() {
	case "pull":
		s = &c.Pull
	case "push":
		s = &c.Push
	}

	if s.Script == "" {
		s = nil
		return
	}

	if s.Shell == "" {
		s.Shell = "bash"
	}

	return
}
