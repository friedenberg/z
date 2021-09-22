package kasten

import "strconv"

type RemoteFileHandle struct {
	Unix int64
	Ext  string
}

func (h RemoteFileHandle) FileName() (fn string) {
	fi := strconv.FormatInt(h.Unix, 10)

	if h.Ext == "" {
		fn = fi
	} else {
		fn = fi + "." + h.Ext
	}

	return
}
