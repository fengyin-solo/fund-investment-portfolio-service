package retry

import "context"

func Run(ctx context.Context, attempts int, call func(context.Context) error) error {
	var last error
	for attempt := 0; attempt < attempts; attempt++ {
		last = call(ctx)
		if last == nil {
			return nil
		}
	}
	return last
}
