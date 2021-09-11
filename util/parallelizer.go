package util

import (
	"sync"

	"golang.org/x/xerrors"
)

type ParallelizerIterFunc func(int, string) error
type ParallelizerErrorFunc func(int, string, error)

type Parallelizer struct {
	Args []string
}

func (p Parallelizer) Run(f ParallelizerIterFunc, e ParallelizerErrorFunc) {
	wg := &sync.WaitGroup{}
	runRead := func() {
		for i, file := range p.Args {
			go func(i int, s string) {
				defer wg.Done()

				err := f(i, s)

				if err != nil {
					err = xerrors.Errorf("%s:\n\t%w", s, err)
					e(i, s, err)
				}
			}(i, file)
		}
	}

	wg.Add(len(p.Args))

	go runRead()
	wg.Wait()
}
