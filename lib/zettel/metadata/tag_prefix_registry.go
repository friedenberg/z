package metadata

type tagConstructor func() ITag

var (
	tagPrefixes = map[string]tagConstructor{}
)

func registerTagPrefix(n string, c tagConstructor) {
	if _, ok := tagPrefixes[n]; ok {
		panic("tag prefix added more than once: " + n)
	}

	tagPrefixes[n] = c
}
