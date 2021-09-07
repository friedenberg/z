package kasten

type Implementation interface {
	InitFromOptions(map[string]interface{}) (err error)
}
