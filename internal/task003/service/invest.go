package service

import (
	"fundinvest/internal/task003/policy"
	"fundinvest/internal/task003/records"
	"fundinvest/internal/task003/remote"
)

type Processor struct {
	Gateway *remote.Gateway
	Records *records.Records
}

func (p *Processor) Submit(key, mode string) error {
	tx := p.Records.Begin(key)
	var firstErr error
	for attempt := 0; attempt < 3; attempt++ {
		err := p.Gateway.Invest(mode)
		if err == nil {
			tx.Commit()
			return firstErr
		}
		if firstErr == nil {
			firstErr = err
		}
		if !policy.Retryable(err) {
			tx.Rollback()
			return err
		}
	}
	tx.Rollback()
	return firstErr
}
