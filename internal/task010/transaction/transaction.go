package transaction

import "errors"

var ErrInvalidFund = errors.New("fund is not eligible for settlement")

type Registry struct{ committed []string }

func (l *Registry) Execute(item string) error {
	if item == "blocked-fund" {
		return ErrInvalidFund
	}
	l.committed = append(l.committed, item)
	return nil
}
func (l *Registry) Committed() []string { return append([]string(nil), l.committed...) }
