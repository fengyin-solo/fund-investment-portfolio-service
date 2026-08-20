package flow

import (
	"context"
	"errors"
	"testing"
	"time"

	"fundinvest/internal/task005/coordinator"
	"fundinvest/internal/task005/producer"
)

func TestSourceFailureStopsBatchWithoutPartialResults(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	result, err := coordinator.Import(ctx, []string{"fund-a", "fund-b", "fund-c"}, 1)
	if !errors.Is(err, producer.ErrSource) {
		t.Errorf("source failure was masked: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("failed batch exposed partial results: %#v", result)
	}
}
