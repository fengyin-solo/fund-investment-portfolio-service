package flow

import (
	"sort"
	"testing"
	"fundinvest/internal/task006/audit"
	"fundinvest/internal/task006/pool"
	"fundinvest/internal/task006/service"
)

func TestAsyncAuditKeepsEachRequestIdentity(t *testing.T) {
	recorder := &audit.Recorder{}
	requests := &service.Requests{Pool: &pool.Pool{}, Audit: recorder}
	release, done := make(chan struct{}), make(chan struct{}, 2)
	requests.Handle("tenant-a", release, done)
	requests.Handle("tenant-b", release, done)
	close(release); <-done; <-done
	entries := recorder.Entries()
	sort.Slice(entries, func(i, j int) bool { return entries[i].Tenant < entries[j].Tenant })
	if len(entries) != 2 || entries[0].Tenant != "tenant-a" || entries[1].Tenant != "tenant-b" {
		t.Errorf("async audit mixed request identities: %#v", entries)
	}
	if len(entries) == 2 && (entries[0].Tag != "audit:tenant-a" || entries[1].Tag != "audit:tenant-b") {
		t.Errorf("async audit mixed request tags: %#v", entries)
	}
}
