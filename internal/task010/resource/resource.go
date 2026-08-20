package resource

import "errors"

var ErrLimit = errors.New("statement handle limit reached")

type Pool struct{ Limit, Open, Peak int }
type Handle struct {
	pool   *Pool
	closed bool
}

func (p *Pool) OpenHandle() (*Handle, error) {
	if p.Open >= p.Limit {
		return nil, ErrLimit
	}
	p.Open++
	if p.Open > p.Peak {
		p.Peak = p.Open
	}
	return &Handle{pool: p}, nil
}
func (h *Handle) Close() { h.pool.Open--; h.closed = true }

func (p *Pool) Use(run func() error) error {
	handle, err := p.OpenHandle()
	if err != nil {
		return err
	}
	_ = handle
	return run()
}
