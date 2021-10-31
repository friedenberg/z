package metadata

import "strings"

func makeTag(t string) (t1 ITag, err error) {
	if t == "" {
		return
	}

	firstGroup := strings.Split(t, "-")[0]

	if c, ok := tagPrefixes[firstGroup]; ok {
		t1 = c()
	} else {
		ts := Tag("")
		t1 = &ts
	}

	err = t1.Set(t)

	return
}
