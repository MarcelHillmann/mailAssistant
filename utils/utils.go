package utils

import "io"

func Closer(c io.Closer) {
	Defer(c.Close)
}

func Defer(c func() error) {
	_ = c()
}
