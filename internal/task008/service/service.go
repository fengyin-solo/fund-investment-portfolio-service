package service

import (
	"context"

	"fundinvest/internal/task008/client"
	"fundinvest/internal/task008/dispatch"
	"fundinvest/internal/task008/retry"
	"fundinvest/internal/task008/tracker"
)

type Quotes struct {
	Client *client.Client
	Tasks  *tracker.Group
}

func (q *Quotes) Start(ctx context.Context) {
	requestCtx := dispatch.RequestContext(ctx)
	q.Tasks.Go(func() error {
		return retry.Run(requestCtx, 3, q.Client.Call)
	})
}
