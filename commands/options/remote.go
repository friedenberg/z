package options

import "golang.org/x/xerrors"

const (
	RemoteCommandPull = RemoteCommand(iota)
	RemoteCommandPush
)

type RemoteCommand int8

func (v RemoteCommand) String() string {
	return ""
}

func (v *RemoteCommand) Set(s string) (err error) {
	switch s {
	case "pull":
		*v = RemoteCommandPull
	case "push":
		*v = RemoteCommandPush
	default:
		err = xerrors.Errorf("unsupported command: '%s'", s)
	}

	return
}
