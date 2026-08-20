package flow

import (
	"testing"

	"fundinvest/internal/task003/records"
	"fundinvest/internal/task003/remote"
	"fundinvest/internal/task003/service"
)

func TestRejectedAndTemporaryInvestmentsKeepDistinctOutcomes(t *testing.T) {
	rejected := &service.Processor{Gateway: &remote.Gateway{}, Records: &records.Records{}}
	if err := rejected.Submit("order-reject", "reject"); err == nil {
		t.Error("business rejection was reported as successful")
	}
	if rejected.Gateway.Calls != 1 || rejected.Records.Count() != 0 {
		t.Errorf("rejection retried or left an investment record: calls=%d entries=%d", rejected.Gateway.Calls, rejected.Records.Count())
	}

	temporary := &service.Processor{Gateway: &remote.Gateway{}, Records: &records.Records{}}
	if err := temporary.Submit("order-recover", "temporary"); err != nil {
		t.Errorf("temporary failure still surfaced after recovery: %v", err)
	}
	if temporary.Gateway.Calls != 2 || temporary.Records.Count() != 1 {
		t.Errorf("recovered investment has wrong side effects: calls=%d entries=%d", temporary.Gateway.Calls, temporary.Records.Count())
	}
}
