package publisher

import "errors"

var ErrTemporary = errors.New("event broker unavailable")

type Bus struct {
	Calls, Delivered int
	accepted         map[string]bool
}

func (b *Bus) Publish(key string) error {
	b.Calls++
	b.Delivered++
	if b.Calls == 1 {
		return ErrTemporary
	}
	return nil
}
