package collector

import "context"

func Drain(ctx context.Context, values <-chan string, errs <-chan error) ([]string, error) {
	var result []string
	for {
		select {
		case value, ok := <-values:
			if !ok { return result, nil }
			result = append(result, value)
		case <-ctx.Done():
			return result, ctx.Err()
		}
	}
}
