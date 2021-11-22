package beginner

import "io"

type Beginner interface {
	Begin(io.Writer)
}
