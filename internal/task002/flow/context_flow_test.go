package flow

import (
	"context"
	"testing"

	"fundinvest/internal/task002/adapter"
	"fundinvest/internal/task002/service"
)

func TestCanceledQuoteDoesNotRetryOrPoisonLaterRequest(t *testing.T) {
	pricing := &adapter.Pricing{}
	quotes := service.New(pricing)
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := quotes.Current(canceled, "F001"); err == nil {
		t.Error("canceled quote unexpectedly succeeded")
	}
	if calls := pricing.Calls(); calls != 0 {
		t.Errorf("canceled quote reached pricing %d times", calls)
	}
	value, err := quotes.Current(context.Background(), "F001")
	if err != nil || value != 12345 {
		t.Errorf("fresh request inherited the earlier cancellation: value=%d err=%v", value, err)
	}
}
