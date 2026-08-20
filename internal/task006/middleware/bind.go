package middleware

import "fundinvest/internal/task006/identity"

func Bind(value *identity.Lease, tenant string) {
	value.Tenant = tenant
	value.Tags = append(value.Tags, "audit:"+tenant)
}
