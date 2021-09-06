package lib

import "github.com/friedenberg/z/lib/kasten"

type Umwelt struct {
	DefaultKasten kasten.Implementation
	Kasten        map[string]kasten.Implementation
}
