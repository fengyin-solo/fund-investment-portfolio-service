package flow

import (
	"fundinvest/internal/task009/audit"
	"fundinvest/internal/task009/publisher"
	"fundinvest/internal/task009/repository"
	"fundinvest/internal/task009/service"
	"testing"
)

func TestPublishRetryKeepsOneCommitEventAndFinalAudit(t *testing.T) {
	repo, bus, log := &repository.Repository{}, &publisher.Bus{}, &audit.Log{}
	settlement := &service.Settlement{Repo: repo, Bus: bus, Audit: log}
	if err := settlement.Run("order-9"); err != nil {
		t.Fatalf("settlement did not recover: %v", err)
	}
	if repo.Begins != 1 || repo.Writes() != 1 {
		t.Errorf("transaction was repeated or duplicated: begins=%d writes=%d", repo.Begins, repo.Writes())
	}
	if bus.Calls != 2 || bus.Delivered != 1 {
		t.Errorf("published event was duplicated: calls=%d delivered=%d", bus.Calls, bus.Delivered)
	}
	states := log.States()
	if len(states) != 1 || states[0] != "success" {
		t.Errorf("audit retained intermediate states: %#v", states)
	}
}
