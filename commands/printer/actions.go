package printer

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
	ActionOpenAll = ActionEdit | ActionOpenFile | ActionOpenUrl
)

func (a *Actions) String() string {
	//TODO
	return ""
	// sb := &strings.Builder{}
}

func (a *Actions) Set(s string) (err error) {
	if s == "" {
		*a = ActionEdit
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
		case "open-all":
			*a = *a | ActionOpenAll
		}
	}

	return
}

func (a *Actions) MatchZettel(z *lib.Zettel) bool {
	if *a&ActionEdit != 0 {
		return true
	}

	if *a&ActionOpenFile != 0 && z.HasFile() {
		return true
	}

	if *a&ActionOpenUrl != 0 && z.HasUrl() {
		return true
	}

	return false
}
