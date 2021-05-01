package concurrency

import "errors"

var (
	ErrConcurrent = errors.New("concurrent update")
)
