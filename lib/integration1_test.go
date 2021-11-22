package lib

import (
	"testing"
)

func TestAddZettelThenReadIndex(t *testing.T) {
	u := makeUmwelt(t)

	contents := `---
- the title
- a-tag
...
`

	z := makeZettel(t, u, contents)
	assertTransactionSuccessful(t, u, 1)
	assertZettelMatchesFileSystem(t, z, contents)
	assertZettelExistsInIndex(t, u, z)
	assertTagInIndex(t, u, z, "a-tag")
}

func TestAddTwoZettelsThenReadIndex(t *testing.T) {
	u := makeUmwelt(t)

	c1 := `---
- the title
- a-tag
...
`

	c2 := `---
- the title 2
- a-tag
- zz-inbox
...
`

	z1 := makeZettel(t, u, c1)
	z2 := makeZettel(t, u, c2)

	assertTransactionSuccessful(t, u, 2)

	assertZettelMatchesFileSystem(t, z1, c1)
	assertZettelExistsInIndex(t, u, z1)
	assertTagInIndex(t, u, z1, "a-tag")

	assertZettelMatchesFileSystem(t, z2, c2)
	assertZettelExistsInIndex(t, u, z2)
	assertTagInIndex(t, u, z2, "a-tag")
	assertTagInIndex(t, u, z2, "zz-inbox")
}

func TestAddTwoZettelsThenModifyOne(t *testing.T) {
	u := makeUmwelt(t)

	c1 := `---
- the title
- a-tag
...
`

	c2 := `---
- the title 2
- a-tag
- zz-inbox
...
`

	z1 := makeZettel(t, u, c1)
	z2 := makeZettel(t, u, c2)

	assertTransactionSuccessful(t, u, 2)

	u.Transaction = MakeTransaction()

	z2.Metadata.SetStringTags([]string{"the title 2"})
	z2.Write(nil)
	u.Set(z2, TransactionActionModified)

	assertTransactionSuccessful(t, u, 1)

	assertZettelMatchesFileSystem(t, z1, c1)
	assertZettelExistsInIndex(t, u, z1)
	assertTagInIndex(t, u, z1, "a-tag")

	assertZettelMatchesFileSystem(t, z2,
		`---
- the title 2
...
`,
	)
	assertZettelExistsInIndex(t, u, z2)
	assertNoTags(t, u, z2, "a-tag")
	assertNoTags(t, u, z2, "zz-inbox")
}
