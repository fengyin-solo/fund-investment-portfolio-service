package identity

type Lease struct { Tenant string; Tags []string }

func Clone(value *Lease) *Lease { return value }
