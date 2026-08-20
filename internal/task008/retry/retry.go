package retry

import "context"

// Run calls call up to attempts times, stopping early on success. It honors
// ctx cancellation: before every attempt it checks ctx, and once canceled it
// returns ctx.Err() without dispatching further work. This is what prevents
// "residual" retry attempts from being launched after a cancel.
func Run(ctx context.Context, attempts int, call func(context.Context) error) error {
	var last error
	for attempt := 0; attempt < attempts; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		last = call(ctx)
		if last == nil {
			return nil
		}
	}
	return last
}
