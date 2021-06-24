package commands

import "flag"

func GetSubcommandAutocomplete(f *flag.FlagSet) CommandRunFunc {
	return func(e Env) error {
		return nil
	}
}
