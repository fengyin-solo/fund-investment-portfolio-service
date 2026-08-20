package producer

import "errors"

var ErrSource = errors.New("fund source failed")

func Start(values []string, failAt int) (<-chan string, <-chan error) {
	out := make(chan string)
	errs := make(chan error)
	go func() {
		for index, value := range values {
			if index == failAt {
				errs <- ErrSource
				return
			}
			out <- value
		}
		close(out)
	}()
	return out, errs
}
