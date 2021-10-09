package metadata

import "net/url"

type Url url.URL

func (u *Url) Set(t string) (err error) {
	a, err := url.Parse(t)

	if err != nil {
		return
	}

	*u = Url(*a)

	return
}

func (u Url) Tag() string {
	a := url.URL(u)
	return "u-" + a.String()
}
