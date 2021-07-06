package commands

import (
	"flag"
	"os"
	"os/exec"

	"github.com/friedenberg/z/lib"
)

func GetSubcommandEdit(f *flag.FlagSet) CommandRunFunc {
	var shouldEdit bool
	var shouldOpen bool

	f.BoolVar(&shouldEdit, "edit", true, "")
	f.BoolVar(&shouldOpen, "open", false, "")

	return func(e *lib.Env) (err error) {
		editor, args := getEditor()
		args = append(args, f.Arg(0))

		cmd := exec.Command(editor, args...)
		cmd.Dir = e.BasePath

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout

		cmd.Run()

		//TODO persist what is necessary
		return
	}
}

func getEditor() (e string, a []string) {
	// var ok bool

	// if e, ok = os.LookupEnv("EDITOR"); ok {
	// 	return
	// }

	// if e, ok = os.LookupEnv("VISUAL"); ok {
	// 	return
	// }

	return "open", a
}
