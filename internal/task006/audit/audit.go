package audit

import (
	"sync"
	"fundinvest/internal/task006/identity"
)

type Entry struct { Tenant, Tag string }
type Recorder struct { mu sync.Mutex; entries []Entry }

func (r *Recorder) Later(value *identity.Lease, release <-chan struct{}, done chan<- struct{}) {
	go func() {
		defer func() { done <- struct{}{} }()
		<-release
		entry := Entry{Tenant: value.Tenant}
		if len(value.Tags) > 0 { entry.Tag = value.Tags[0] }
		r.mu.Lock(); r.entries = append(r.entries, entry); r.mu.Unlock()
	}()
}

func (r *Recorder) Entries() []Entry {
	r.mu.Lock(); defer r.mu.Unlock()
	return append([]Entry(nil), r.entries...)
}
