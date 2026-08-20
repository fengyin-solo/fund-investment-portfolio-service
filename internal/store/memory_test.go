package store

import (
	"testing"
	"time"

	"fundinvest/internal/model"
)

func TestMemoryStoreFundCRUD(t *testing.T) {
	s := NewMemoryStore()
	f := &model.Fund{ID: "f1", Code: "000001", Name: "Test Fund", Type: model.FundTypeStock, RiskLevel: 3, Nav: 15000, Status: model.FundStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateFund(f); err != nil {
		t.Fatalf("create fund: %v", err)
	}
	if err := s.CreateFund(&model.Fund{ID: "f2", Code: "000001", Name: "Dup", Type: model.FundTypeStock, RiskLevel: 3, Nav: 100, Status: model.FundStatusActive}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetFund("f1")
	if err != nil {
		t.Fatalf("get fund: %v", err)
	}
	if got.ID != "f1" {
		t.Fatalf("unexpected id %s", got.ID)
	}
	if _, err := s.GetFund("notexist"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	f.Name = "Updated"
	if err := s.UpdateFund(f); err != nil {
		t.Fatalf("update fund: %v", err)
	}
	if err := s.DeleteFund("f1"); err != nil {
		t.Fatalf("delete fund: %v", err)
	}
	if err := s.DeleteFund("f1"); err != ErrNotFound {
		t.Fatalf("expected not found on second delete, got %v", err)
	}
}

func TestMemoryStoreAccountCRUD(t *testing.T) {
	s := NewMemoryStore()
	a := &model.Account{ID: "a1", Owner: "Alice", Balance: 10000, CreatedAt: time.Now()}
	if err := s.CreateAccount(a); err != nil {
		t.Fatalf("create account: %v", err)
	}
	got, err := s.GetAccount("a1")
	if err != nil {
		t.Fatalf("get account: %v", err)
	}
	if got.Owner != "Alice" {
		t.Fatalf("unexpected owner %s", got.Owner)
	}
	if _, err := s.GetAccount("notexist"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	a.Balance = 20000
	if err := s.UpdateAccount(a); err != nil {
		t.Fatalf("update account: %v", err)
	}
	if err := s.DeleteAccount("a1"); err != nil {
		t.Fatalf("delete account: %v", err)
	}
	if err := s.DeleteAccount("a1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStorePlanCRUD(t *testing.T) {
	s := NewMemoryStore()
	p := &model.Plan{ID: "p1", AccountID: "a1", FundID: "f1", Amount: 1000, Interval: model.PlanIntervalMonthly, Status: model.PlanStatusActive, NextExecAt: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreatePlan(p); err != nil {
		t.Fatalf("create plan: %v", err)
	}
	got, err := s.GetPlan("p1")
	if err != nil {
		t.Fatalf("get plan: %v", err)
	}
	if got.ID != "p1" {
		t.Fatalf("unexpected id %s", got.ID)
	}
	if _, err := s.GetPlan("notexist"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	p.Amount = 2000
	if err := s.UpdatePlan(p); err != nil {
		t.Fatalf("update plan: %v", err)
	}
	if err := s.DeletePlan("p1"); err != nil {
		t.Fatalf("delete plan: %v", err)
	}
	if err := s.DeletePlan("p1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStoreTransactionCRUD(t *testing.T) {
	s := NewMemoryStore()
	tx := &model.Transaction{ID: "t1", AccountID: "a1", FundID: "f1", Type: model.TransactionTypeInvest, Amount: 1000, Shares: 500, Price: 200, Status: model.TransactionStatusConfirmed, CreatedAt: time.Now()}
	if err := s.CreateTransaction(tx); err != nil {
		t.Fatalf("create transaction: %v", err)
	}
	got, err := s.GetTransaction("t1")
	if err != nil {
		t.Fatalf("get transaction: %v", err)
	}
	if got.ID != "t1" {
		t.Fatalf("unexpected id %s", got.ID)
	}
	if _, err := s.GetTransaction("notexist"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	tx.Amount = 2000
	if err := s.UpdateTransaction(tx); err != nil {
		t.Fatalf("update transaction: %v", err)
	}
	if err := s.DeleteTransaction("t1"); err != nil {
		t.Fatalf("delete transaction: %v", err)
	}
	if err := s.DeleteTransaction("t1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStoreHoldingCRUD(t *testing.T) {
	s := NewMemoryStore()
	h := &model.Holding{ID: "h1", AccountID: "a1", FundID: "f1", TotalShares: 1000, Cost: 5000, MarketValue: 6000, UpdatedAt: time.Now()}
	if err := s.CreateHolding(h); err != nil {
		t.Fatalf("create holding: %v", err)
	}
	got, err := s.GetHolding("h1")
	if err != nil {
		t.Fatalf("get holding: %v", err)
	}
	if got.ID != "h1" {
		t.Fatalf("unexpected id %s", got.ID)
	}
	gh, err := s.GetHoldingByAccountAndFund("a1", "f1")
	if err != nil {
		t.Fatalf("get holding by account and fund: %v", err)
	}
	if gh.ID != "h1" {
		t.Fatalf("unexpected id %s", gh.ID)
	}
	if _, err := s.GetHoldingByAccountAndFund("a2", "f2"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	h.TotalShares = 2000
	if err := s.UpdateHolding(h); err != nil {
		t.Fatalf("update holding: %v", err)
	}
	if err := s.DeleteHolding("h1"); err != nil {
		t.Fatalf("delete holding: %v", err)
	}
	if err := s.DeleteHolding("h1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStoreNavHistoryCRUD(t *testing.T) {
	s := NewMemoryStore()
	n := &model.NavHistory{ID: "n1", FundID: "f1", Nav: 15000, RecordedAt: time.Now()}
	if err := s.CreateNavHistory(n); err != nil {
		t.Fatalf("create nav history: %v", err)
	}
	got, err := s.GetNavHistory("n1")
	if err != nil {
		t.Fatalf("get nav history: %v", err)
	}
	if got.ID != "n1" {
		t.Fatalf("unexpected id %s", got.ID)
	}
	if _, err := s.GetNavHistory("notexist"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	n.Nav = 16000
	if err := s.UpdateNavHistory(n); err != nil {
		t.Fatalf("update nav history: %v", err)
	}
	if err := s.DeleteNavHistory("n1"); err != nil {
		t.Fatalf("delete nav history: %v", err)
	}
	if err := s.DeleteNavHistory("n1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStoreDividendCRUD(t *testing.T) {
	s := NewMemoryStore()
	d := &model.Dividend{ID: "d1", AccountID: "a1", FundID: "f1", Amount: 1000, Type: model.DividendTypeCash, CreatedAt: time.Now()}
	if err := s.CreateDividend(d); err != nil {
		t.Fatalf("create dividend: %v", err)
	}
	got, err := s.GetDividend("d1")
	if err != nil {
		t.Fatalf("get dividend: %v", err)
	}
	if got.ID != "d1" {
		t.Fatalf("unexpected id %s", got.ID)
	}
	if _, err := s.GetDividend("notexist"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	d.Amount = 2000
	if err := s.UpdateDividend(d); err != nil {
		t.Fatalf("update dividend: %v", err)
	}
	list := s.ListDividends()
	if len(list) != 1 {
		t.Fatalf("expected 1 dividend, got %d", len(list))
	}
	if err := s.DeleteDividend("d1"); err != nil {
		t.Fatalf("delete dividend: %v", err)
	}
	if err := s.DeleteDividend("d1"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMemoryStorePointerConsistency(t *testing.T) {
	s := NewMemoryStore()
	f := &model.Fund{ID: "f1", Code: "000100", Name: "Pointer Test", Type: model.FundTypeStock, RiskLevel: 3, Nav: 10000, Status: model.FundStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	s.CreateFund(f)
	f.Name = "Mutated"
	got, _ := s.GetFund("f1")
	if got.Name != "Mutated" {
		t.Fatalf("expected mutated name due to pointer storage")
	}
}

func TestMemoryStoreListFundsReturnsAll(t *testing.T) {
	s := NewMemoryStore()
	s.CreateFund(&model.Fund{ID: "f1", Code: "000101", Name: "A", Type: model.FundTypeStock, RiskLevel: 3, Nav: 100, Status: model.FundStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	s.CreateFund(&model.Fund{ID: "f2", Code: "000102", Name: "B", Type: model.FundTypeBond, RiskLevel: 2, Nav: 200, Status: model.FundStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	list := s.ListFunds()
	if len(list) != 2 {
		t.Fatalf("expected 2 funds, got %d", len(list))
	}
}

func TestMemoryStoreListTransactionsEmpty(t *testing.T) {
	s := NewMemoryStore()
	list := s.ListTransactions()
	if len(list) != 0 {
		t.Fatalf("expected empty list, got %d", len(list))
	}
}

func TestMemoryStoreUpdateNotExist(t *testing.T) {
	s := NewMemoryStore()
	if err := s.UpdateFund(&model.Fund{ID: "none", Code: "000103", Name: "X", Type: model.FundTypeStock, RiskLevel: 3, Nav: 100, Status: model.FundStatusActive}); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

