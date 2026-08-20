package client

import (
	"context"
	"errors"
	"sync/atomic"
)

var ErrTemporary = errors.New("quote endpoint temporarily unavailable")

type Client struct {
	Permit  <-chan struct{}
	Started chan<- int
	calls   atomic.Int64
}

func (c *Client) Call(ctx context.Context) error {
	call := int(c.calls.Add(1))
	c.Started <- call
	// An in-flight quote request must yield to cancellation so that shutting
	// down the caller does not block until the (possibly never-arriving)
	// permit/permit-or-timeout is reached.
	select {
	case <-c.Permit:
		return ErrTemporary
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Client) Calls() int { return int(c.calls.Load()) }
