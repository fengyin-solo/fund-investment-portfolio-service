package service

import (
	"fundinvest/internal/task010/audit"
	"fundinvest/internal/task010/resource"
	"fundinvest/internal/task010/transaction"
)

type Batch struct {
	Pool     *resource.Pool
	Registry *transaction.Registry
	Audit    *audit.Log
}

func (b *Batch) Process(items []string) error {
	for _, item := range items {
		handle, err := b.Pool.OpenHandle()
		if err != nil {
			return err
		}
		defer handle.Close()
		b.Audit.RecordStarted(item)
		if err := b.Registry.Execute(item); err != nil {
			return err
		}
	}
	return nil
}
