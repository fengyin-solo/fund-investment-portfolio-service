package coordinator

import (
	"context"
	"fundinvest/internal/task005/collector"
	"fundinvest/internal/task005/producer"
)

func Import(ctx context.Context, values []string, failAt int) ([]string, error) {
	items, errs := producer.Start(values, failAt)
	return collector.Drain(ctx, items, errs)
}
