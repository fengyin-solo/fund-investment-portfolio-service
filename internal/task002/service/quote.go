package service

import (
	"context"

	"fundinvest/internal/task002/adapter"
	"fundinvest/internal/task002/scope"
	"fundinvest/internal/task002/worker"
)

type Quotes struct {
	router  *scope.Router
	pricing *adapter.Pricing
}

func New(pricing *adapter.Pricing) *Quotes {
	return &Quotes{router: scope.New(), pricing: pricing}
}

func (q *Quotes) Current(ctx context.Context, fund string) (int64, error) {
	requestCtx := q.router.ForRequest(ctx)
	return worker.Do(requestCtx, 3, func(callCtx context.Context) (int64, error) {
		return q.pricing.Fetch(callCtx, fund)
	})
}
