package flow

import (
	"testing"

	"fundinvest/internal/task004/cache"
	"fundinvest/internal/task004/exporter"
	"fundinvest/internal/task004/service"
)

func TestEarlierImportRemainsStableAfterNextBatch(t *testing.T) {
	store, queue := cache.New(), &exporter.Queue{}
	imports := service.New(store, queue)
	imports.Accept("batch-a", "ALPHA")
	imports.Accept("batch-b", "BRAVO")
	first, ok := store.Load("batch-a")
	if !ok || string(first.Payload) != "ALPHA" {
		t.Errorf("cached first batch changed after second import: %#v", first)
	}
	exported := queue.Flush()
	if len(exported) != 2 || exported[0] != "ALPHA" || exported[1] != "BRAVO" {
		t.Errorf("export mixed the two import payloads: %#v", exported)
	}
}
