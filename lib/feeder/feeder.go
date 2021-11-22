package feeder

type feeder struct {
	runFunc func(chan<- string)
	argChan chan string
}

func MakeFeeder(runFunc func(chan<- string)) feeder {
	return feeder{
		runFunc: runFunc,
		argChan: make(chan string, 0),
	}
}

func (f feeder) GetChan() <-chan string {
	return f.argChan
}

func (f feeder) Run() {
	defer close(f.argChan)

	if f.runFunc != nil {
		f.runFunc(f.argChan)
	}
}
