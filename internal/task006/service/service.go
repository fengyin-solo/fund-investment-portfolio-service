package service

import (
	"fundinvest/internal/task006/audit"
	"fundinvest/internal/task006/identity"
	"fundinvest/internal/task006/middleware"
	"fundinvest/internal/task006/pool"
)

type Requests struct { Pool *pool.Pool; Audit *audit.Recorder }

func (r *Requests) Handle(tenant string, release <-chan struct{}, done chan<- struct{}) {
	lease := r.Pool.Acquire()
	middleware.Bind(lease, tenant)
	r.Audit.Later(identity.Clone(lease), release, done)
	r.Pool.Release(lease)
}
