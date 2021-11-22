package feeder

type StringSlice struct {
	args    []string
	argChan chan string
}

func MakeStringSlice(args []string) StringSlice {
	return StringSlice{
		args:    args,
		argChan: make(chan string, 0),
	}
}

func (f StringSlice) GetChan() <-chan string {
	return f.argChan
}

func (f StringSlice) Run() {
	defer close(f.argChan)

	for _, a := range f.args {
		f.argChan <- a
	}
}
