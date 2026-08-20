package dispatch

import "context"

func RequestContext(ctx context.Context) context.Context {
	return context.Background()
}
