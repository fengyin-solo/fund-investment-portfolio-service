package builder

import "fmt"

type Portfolio struct {
	ID        string
	Positions map[string]int64
	Ready     bool
}

type Builder struct{}

func (Builder) Build(id string, interrupt bool) (result *Portfolio, err error) {
	result = &Portfolio{ID: id, Positions: make(map[string]int64)}
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("portfolio build interrupted: %v", recovered)
		}
	}()

	result.Positions["bond"] = 6000
	if interrupt {
		panic("price source unavailable")
	}
	result.Positions["equity"] = 4000
	result.Ready = true
	return result, nil
}
