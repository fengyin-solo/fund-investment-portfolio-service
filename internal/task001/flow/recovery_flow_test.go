package flow

import (
	"testing"

	"fundinvest/internal/task001/cache"
	"fundinvest/internal/task001/service"
)

func TestInterruptedPortfolioDoesNotLeakIntoNextLookup(t *testing.T) {
	manager := service.New(cache.New())
	if err := manager.Prepare("client-a", true); err == nil {
		t.Fatal("interrupted preparation should report an error")
	}
	if leaked, ok := manager.Lookup("client-a"); ok {
		t.Fatalf("failed preparation leaked a partial portfolio: %#v", leaked)
	}
	if err := manager.Prepare("client-b", false); err != nil {
		t.Fatalf("later preparation failed: %v", err)
	}
	portfolio, ok := manager.Lookup("client-b")
	if !ok || !portfolio.Ready || len(portfolio.Positions) != 2 {
		t.Fatalf("later lookup did not return an isolated complete portfolio: %#v", portfolio)
	}
}
