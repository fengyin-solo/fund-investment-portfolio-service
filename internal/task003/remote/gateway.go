package remote

import (
	"errors"
	"fmt"
)

var ErrRejected = errors.New("investment rejected")

type TemporaryError struct{ Message string }

func (e *TemporaryError) Error() string { return e.Message }

type Gateway struct{ Calls int }

func (g *Gateway) Invest(mode string) error {
	g.Calls++
	var err error
	switch mode {
	case "reject":
		err = ErrRejected
	case "temporary":
		if g.Calls == 1 {
			err = &TemporaryError{Message: "pricing service busy"}
		}
	}
	if err != nil {
		return fmt.Errorf("investment unavailable: %v", err)
	}
	return nil
}
