package kasten

import "fmt"

type TRegistry map[string]Implementation

var (
	Registry TRegistry
)

func init() {
	Registry = TRegistry(make(map[string]Implementation))
}

func (r TRegistry) Register(n string, i Implementation) (err error) {
	if _, ok := map[string]Implementation(r)[n]; ok {
		err = fmt.Errorf("Multiple implementations with name: '%s'", n)
		return
	}

	map[string]Implementation(r)[n] = i

	return
}

func (r TRegistry) Get(n string) (i Implementation, ok bool) {
	i, ok = map[string]Implementation(r)[n]
	return
}
