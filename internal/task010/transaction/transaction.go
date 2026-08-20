package transaction

import "errors"

var ErrInvalidFund = errors.New("fund is not eligible for settlement")

type Registry struct{ committed []string }

func (l *Registry) Execute(item string) (err error) {
	defer func() { l.committed = append(l.committed, item); err = nil }()
	if item == "blocked-fund" {
		return ErrInvalidFund
	}
	return nil
}
func (l *Registry) Committed() []string { return append([]string(nil), l.committed...) }
