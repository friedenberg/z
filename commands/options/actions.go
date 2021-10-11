package options

import (
	"strings"

	"github.com/friedenberg/z/lib"
)

type Actions uint8

const (
	ActionUnknown = 0
	ActionEdit    = 1 << iota
	ActionOpenFile
	ActionOpenUrl
	ActionPrintZettelPath
	ActionOpenAll = ActionEdit | ActionOpenFile | ActionOpenUrl | ActionPrintZettelPath
)

func (a *Actions) String() string {
	//TODO
	return ""
	// sb := &strings.Builder{}
}

func (a *Actions) Set(s string) (err error) {
	if *a == ActionEdit {
		*a = ActionUnknown
	}

	if s == "" {
		*a = *a | ActionEdit
		return
	}

	actions := strings.Split(s, ",")

	for _, action := range actions {
		switch action {
		case "edit":
			*a = *a | ActionEdit
		case "open-files":
			*a = *a | ActionOpenFile
		case "open-urls":
			*a = *a | ActionOpenUrl
		case "print-zettel-path":
			*a = *a | ActionPrintZettelPath
		case "open-all":
			*a = *a | ActionOpenAll
		}
	}

	return
}

func (a *Actions) ShouldEdit() bool {
	return *a&ActionEdit != 0
}

func (a *Actions) ShouldOpenFile() bool {
	return *a&ActionOpenFile != 0
}

func (a *Actions) ShouldOpenUrl() bool {
	return *a&ActionOpenUrl != 0
}

func (a *Actions) MatchZettel(z *lib.Zettel) bool {
	if *a&ActionEdit != 0 {
		return true
	}

	if *a&ActionOpenFile != 0 && z.Note.Metadata.HasFile() {
		return true
	}

	_, hu := z.Note.Metadata.Url()

	if *a&ActionOpenUrl != 0 && hu {
		return true
	}

	if *a&ActionPrintZettelPath != 0 {
		return true
	}

	return false
}
