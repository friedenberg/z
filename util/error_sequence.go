package util

type ErrorSequin func() error

type ErrorSequence struct {
	Sequins []ErrorSequin
}

func (s ErrorSequence) Run() (err error) {
	for _, e := range s.Sequins {
		err = e()

		if err != nil {
			return
		}
	}

	return
}
