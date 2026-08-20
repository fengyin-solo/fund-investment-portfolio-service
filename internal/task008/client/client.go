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
	<-c.Permit
	return ErrTemporary
}

func (c *Client) Calls() int { return int(c.calls.Load()) }
