package worker

import "context"

func Do(ctx context.Context, attempts int, fn func(context.Context) (int64, error)) (int64, error) {
	var last error
	for i := 0; i < attempts; i++ {
		value, err := fn(ctx)
		if err == nil {
			return value, nil
		}
		last = err
	}
	return 0, last
}
