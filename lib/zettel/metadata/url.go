package metadata

import (
	"net/url"

	"github.com/friedenberg/z/util"
	"golang.org/x/xerrors"
)

type Url url.URL

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

	*u = Url(*a)

	return
}

func (u Url) Tag() string {
	return "u-" + u.String()
}

func (u Url) String() string {
	a := url.URL(u)
	b := a.String()

	return b
}

func (u Url) SearchMatchTags() (expanded TagSet) {
	expanded = MakeTagSet()
	expanded.Merge(Tag("d-" + u.Host).SearchMatchTags())

	return
}
