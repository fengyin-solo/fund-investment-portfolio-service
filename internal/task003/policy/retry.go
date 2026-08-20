package policy

func Retryable(err error) bool { return err != nil }
