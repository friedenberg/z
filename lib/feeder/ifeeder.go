package feeder

type Feeder interface {
	GetChan() <-chan string
	Run()
}
