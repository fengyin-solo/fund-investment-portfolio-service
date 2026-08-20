package publisher

import "errors"

var ErrTemporary = errors.New("event broker unavailable")

type Bus struct {
	Calls, Delivered int
	accepted         map[string]bool
}

func (b *Bus) Publish(key string) error {
	b.Calls++
	if b.accepted == nil {
		b.accepted = make(map[string]bool)
	}
	// A retried publish of the same event must not be delivered twice.
	if b.accepted[key] {
		return nil
	}
	if b.Calls == 1 {
		return ErrTemporary
	}
	b.accepted[key] = true
	b.Delivered++
	return nil
}
