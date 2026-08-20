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
	if c.seen[key] {
		// Already settled under this idempotency key: dedup the retry.
		return nil
	}
	if c.Calls == 1 {
		// Transient failure settles nothing; the caller retries the same key.
		return ErrTemporary
	}
	c.seen[key] = true
	c.SideEffects++
	return nil
}
