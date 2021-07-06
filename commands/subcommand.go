package commands

import "github.com/friedenberg/z/lib"

type CommandRunFunc func(*lib.Env) error
