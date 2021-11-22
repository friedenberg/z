package metadata

import (
	"net/url"
	"os"
	"strings"

	"github.com/friedenberg/z/util"
	"github.com/friedenberg/z/util/stdprinter"
	"golang.org/x/xerrors"
)

func init() {
	registerTagPrefix(
		"u",
		func() (t ITag) { return &Url{} },
	)
}

type Url struct {
	url.URL
}

func (u *Url) Set(t string) (err error) {
	if len(t) < 3 {
		err = xerrors.Errorf("too few characters for URL")
		return
	}

	if t[0:2] != "u-" {
		err = xerrors.Errorf("missing u- prefix")
		return
	}

	t = t[2:]

	a, err := util.ParseURL(t)

	if err != nil {
		return
	}

	if a.Hostname() == "" {
		err = xerrors.Errorf("hostname for url ('%s') is empty", t)
		return
	}

	u.URL = *a

	return
}

func (u Url) Tag() string {
	return "u-" + u.String()
}

func (u Url) String() string {
	return u.URL.String()
}

func (u Url) CorrectedString() (s string) {
	u1 := u.URL

	s = u1.String()

	if u1.Scheme == "file" && u1.Hostname() == "~" {
		homeDir, err := os.UserHomeDir()
		stdprinter.PanicIfError(err)
		s = strings.Replace(s, "~", homeDir, 1)
	}

	return
}

func (u Url) SearchMatchTags() (expanded TagSet) {
	expanded = MakeTagSet()
	expanded.Merge(Tag("d-" + u.Host).SearchMatchTags())

	return
}
