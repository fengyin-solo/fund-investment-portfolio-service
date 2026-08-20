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
		if err := b.processItem(item); err != nil {
			return err
		}
	}
	return nil
}

// processItem runs one item in its own scope so the handle is released as soon
// as the item finishes, never accumulating across iterations.
func (b *Batch) processItem(item string) (err error) {
	handle, err := b.Pool.OpenHandle()
	if err != nil {
		return err
	}
	defer handle.Close()

	b.Audit.RecordStarted(item)
	if err := b.Registry.Execute(item); err != nil {
		return err
	}
	b.Audit.RecordCommitted(item)
	return nil
}
