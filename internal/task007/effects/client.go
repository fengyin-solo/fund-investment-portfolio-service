package effects

import "errors"

var ErrTemporary = errors.New("temporary settlement failure")

type Client struct {
	seen               map[string]bool
	Calls, SideEffects int
}

func New() *Client { return &Client{seen: map[string]bool{}} }
func (c *Client) Execute(key string) error {
	c.Calls++
	if !c.seen[key] {
		c.SideEffects++
	}
	if c.Calls == 1 {
		return ErrTemporary
	}
	c.seen[key] = true
	return nil
}
