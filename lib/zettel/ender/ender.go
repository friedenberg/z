package ender

import "io"

type Ender interface {
	End(io.Writer)
}
