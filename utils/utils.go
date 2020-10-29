package utils

import "io"

// Closer is helper function
func Closer(c io.Closer) {
	Defer(c.Close)
}

// Defer is helper function
func Defer(c func() error) {
	_ = c()
}
