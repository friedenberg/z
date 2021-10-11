package metadata

import "testing"

func TestUrl(t *testing.T) {
	assertUrlSet(t,
		`u-google.com`,
		`u-https://google.com`,
	)
}

func TestUrl2(t *testing.T) {
	assertUrlSet(t,
		`u-www.google.com`,
		`u-https://www.google.com`,
	)
}

func TestUrl3(t *testing.T) {
	assertUrlSet(t,
		`u-http://www.google.com`,
		`u-https://www.google.com`,
	)
}

func TestUrl4(t *testing.T) {
	assertUrlSet(t,
		`u-https://feministing.com/2014/05/30/an-open-letter-to-privileged-people-who-play-devils-advocate/`,
		`u-https://feministing.com/2014/05/30/an-open-letter-to-privileged-people-who-play-devils-advocate/`,
	)
}

func assertUrlSet(t *testing.T, s, expected string) {
	t.Helper()

	var u Url
	err := u.Set(s)

	if err != nil {
		t.Errorf("failed to set from tag: %w", err)
	}

	actual := u.Tag()

	if expected != actual {
		t.Errorf("Actual tag was '%s', wanted '%s'", actual, expected)
	}
}
