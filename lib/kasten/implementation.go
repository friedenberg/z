package kasten

type Implementation interface {
	InitFromOptions(map[string]string) (err error)
}
