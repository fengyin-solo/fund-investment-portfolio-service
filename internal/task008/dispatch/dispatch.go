package dispatch

import "context"

// RequestContext carries the caller's cancellation into the retry/dispatch
// path. It must not detach from ctx — otherwise a cancel on the caller side
// never reaches the in-flight work, and retries keep being dispatched.
func RequestContext(ctx context.Context) context.Context {
	return ctx
}
