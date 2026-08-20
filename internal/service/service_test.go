package service

import (
	"strconv"
	"testing"
	"time"

	"fundinvest/internal/config"
	"fundinvest/internal/model"
	"fundinvest/internal/store"
	"fundinvest/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func TestCreateFundAndAccount(t *testing.T) {
	svc := newTestService()
	f, err := svc.CreateFund(model.Fund{Code: "000001", Name: "Test Fund", Type: model.FundTypeStock, RiskLevel: 3, Nav: 15000})
	if err != nil {
		t.Fatalf("create fund: %v", err)
	}
	if f.ID == "" {
		t.Fatal("fund id empty")
	}
	a, err := svc.CreateAccount(model.Account{Owner: "Alice", Balance: 50000})
	if err != nil {
		t.Fatalf("create account: %v", err)
	}
	if a.ID == "" {
		t.Fatal("account id empty")
	}
}

func TestCreatePlanCrossEntityValidation(t *testing.T) {
	svc := newTestService()
	_, err := svc.CreatePlan(model.Plan{AccountID: "no", FundID: "no", Amount: 1000, Interval: model.PlanIntervalMonthly})
	if err == nil {
		t.Fatal("expected error for missing account")
	}
	acc, _ := svc.CreateAccount(model.Account{Owner: "Bob", Balance: 10000})
	_, err = svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: "no", Amount: 1000, Interval: model.PlanIntervalMonthly})
	if err == nil {
		t.Fatal("expected error for missing fund")
	}
	fund, _ := svc.CreateFund(model.Fund{Code: "000002", Name: "F", Type: model.FundTypeBond, RiskLevel: 2, Nav: 10000})
	_, err = svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: fund.ID, Amount: 1000, Interval: model.PlanIntervalMonthly})
	if err != nil {
		t.Fatalf("create plan failed: %v", err)
	}
}

func TestExecutePlanSuccess(t *testing.T) {
	svc := newTestService()
	fund, _ := svc.CreateFund(model.Fund{Code: "000003", Name: "F3", Type: model.FundTypeStock, RiskLevel: 3, Nav: 20000})
	acc, _ := svc.CreateAccount(model.Account{Owner: "C", Balance: 100000})
	plan, _ := svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: fund.ID, Amount: 10000, Interval: model.PlanIntervalMonthly, NextExecAt: time.Now().Add(-time.Hour)})

	txn, err := svc.ExecutePlan(plan.ID)
	if err != nil {
		t.Fatalf("execute plan: %v", err)
	}
	if txn.Status != model.TransactionStatusConfirmed {
		t.Fatalf("expected confirmed, got %s", txn.Status)
	}
	if txn.Amount != 10000 {
		t.Fatalf("expected amount 10000, got %d", txn.Amount)
	}
	acc2, _ := svc.GetAccount(acc.ID)
	if acc2.Balance != 90000 {
		t.Fatalf("expected balance 90000, got %d", acc2.Balance)
	}
	holdings, _, _ := svc.ListHoldings(model.HoldingFilter{AccountID: acc.ID}, 1, 10)
	if len(holdings) != 1 {
		t.Fatalf("expected 1 holding, got %d", len(holdings))
	}
	if holdings[0].TotalShares != (10000*10000)/20000 {
		t.Fatalf("unexpected shares %d", holdings[0].TotalShares)
	}
	plan2, _ := svc.GetPlan(plan.ID)
	if plan2.InvestedCount != 1 {
		t.Fatalf("expected invested count 1, got %d", plan2.InvestedCount)
	}
}

func TestExecutePlanInsufficientBalance(t *testing.T) {
	svc := newTestService()
	fund, _ := svc.CreateFund(model.Fund{Code: "000004", Name: "F4", Type: model.FundTypeBond, RiskLevel: 2, Nav: 10000})
	acc, _ := svc.CreateAccount(model.Account{Owner: "D", Balance: 500})
	plan, _ := svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: fund.ID, Amount: 1000, Interval: model.PlanIntervalDaily, NextExecAt: time.Now().Add(-time.Hour)})

	txn, err := svc.ExecutePlan(plan.ID)
	if err != nil {
		t.Fatalf("execute plan should not error: %v", err)
	}
	if txn.Status != model.TransactionStatusFailed {
		t.Fatalf("expected failed, got %s", txn.Status)
	}
	acc2, _ := svc.GetAccount(acc.ID)
	if acc2.Balance != 500 {
		t.Fatalf("balance should remain 500, got %d", acc2.Balance)
	}
}

func TestRedeem(t *testing.T) {
	svc := newTestService()
	fund, _ := svc.CreateFund(model.Fund{Code: "000005", Name: "F5", Type: model.FundTypeMixed, RiskLevel: 3, Nav: 10000})
	acc, _ := svc.CreateAccount(model.Account{Owner: "E", Balance: 50000})
	plan, _ := svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: fund.ID, Amount: 20000, Interval: model.PlanIntervalWeekly, NextExecAt: time.Now().Add(-time.Hour)})
	svc.ExecutePlan(plan.ID)

	txn, err := svc.Redeem(acc.ID, fund.ID, 5000)
	if err != nil {
		t.Fatalf("redeem: %v", err)
	}
	if txn.Status != model.TransactionStatusConfirmed {
		t.Fatalf("expected confirmed, got %s", txn.Status)
	}
	if txn.Type != model.TransactionTypeRedeem {
		t.Fatalf("expected redeem type")
	}
	acc2, _ := svc.GetAccount(acc.ID)
	expectedBalance := int64(35000)
	if acc2.Balance != expectedBalance {
		t.Fatalf("expected balance %d, got %d", expectedBalance, acc2.Balance)
	}
	holdings, _, _ := svc.ListHoldings(model.HoldingFilter{AccountID: acc.ID}, 1, 10)
	if len(holdings) != 1 {
		t.Fatalf("expected 1 holding, got %d", len(holdings))
	}
	if holdings[0].TotalShares != (20000*10000)/10000-5000 {
		t.Fatalf("unexpected shares %d", holdings[0].TotalShares)
	}
}

func TestRedeemInsufficientShares(t *testing.T) {
	svc := newTestService()
	fund, _ := svc.CreateFund(model.Fund{Code: "000006", Name: "F6", Type: model.FundTypeIndex, RiskLevel: 2, Nav: 10000})
	acc, _ := svc.CreateAccount(model.Account{Owner: "F", Balance: 0})
	plan, _ := svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: fund.ID, Amount: 10000, Interval: model.PlanIntervalMonthly, NextExecAt: time.Now().Add(-time.Hour)})
	svc.ExecutePlan(plan.ID)

	_, err := svc.Redeem(acc.ID, fund.ID, 20000)
	if err == nil {
		t.Fatal("expected error for insufficient shares")
	}
}

func TestPlanStateMachine(t *testing.T) {
	svc := newTestService()
	fund, _ := svc.CreateFund(model.Fund{Code: "000007", Name: "F7", Type: model.FundTypeStock, RiskLevel: 3, Nav: 10000})
	acc, _ := svc.CreateAccount(model.Account{Owner: "G", Balance: 10000})
	plan, _ := svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: fund.ID, Amount: 1000, Interval: model.PlanIntervalMonthly})

	_, err := svc.ChangePlanStatus(plan.ID, model.PlanStatusPaused)
	if err != nil {
		t.Fatalf("active->paused: %v", err)
	}
	_, err = svc.ChangePlanStatus(plan.ID, model.PlanStatusActive)
	if err != nil {
		t.Fatalf("paused->active: %v", err)
	}
	_, err = svc.ChangePlanStatus(plan.ID, model.PlanStatusClosed)
	if err != nil {
		t.Fatalf("active->closed: %v", err)
	}
	_, err = svc.ChangePlanStatus(plan.ID, model.PlanStatusActive)
	if err == nil {
		t.Fatal("closed->active should fail")
	}
}

func TestTransactionFilterAndList(t *testing.T) {
	svc := newTestService()
	fund, _ := svc.CreateFund(model.Fund{Code: "000008", Name: "F8", Type: model.FundTypeBond, RiskLevel: 2, Nav: 10000})
	acc, _ := svc.CreateAccount(model.Account{Owner: "H", Balance: 50000})
	plan, _ := svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: fund.ID, Amount: 5000, Interval: model.PlanIntervalDaily, NextExecAt: time.Now().Add(-time.Hour)})
	svc.ExecutePlan(plan.ID)

	items, total, err := svc.ListTransactions(model.TransactionFilter{AccountID: acc.ID, Type: model.TransactionTypeInvest, Status: model.TransactionStatusConfirmed}, 1, 10)
	if err != nil {
		t.Fatalf("list transactions: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
}

func TestGetAccountProfit(t *testing.T) {
	svc := newTestService()
	fund, _ := svc.CreateFund(model.Fund{Code: "000009", Name: "F9", Type: model.FundTypeStock, RiskLevel: 3, Nav: 10000})
	acc, _ := svc.CreateAccount(model.Account{Owner: "I", Balance: 50000})
	plan, _ := svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: fund.ID, Amount: 10000, Interval: model.PlanIntervalMonthly, NextExecAt: time.Now().Add(-time.Hour)})
	svc.ExecutePlan(plan.ID)

	profit, err := svc.GetAccountProfit(acc.ID)
	if err != nil {
		t.Fatalf("get account profit: %v", err)
	}
	if profit.TotalCost != 10000 {
		t.Fatalf("expected total cost 10000, got %d", profit.TotalCost)
	}
	if profit.TotalMarket != 10000 {
		t.Fatalf("expected total market 10000, got %d", profit.TotalMarket)
	}
	if profit.TotalProfit != 0 {
		t.Fatalf("expected total profit 0, got %d", profit.TotalProfit)
	}
}

func TestGetTopFunds(t *testing.T) {
	svc := newTestService()
	f1, _ := svc.CreateFund(model.Fund{Code: "000010", Name: "F10", Type: model.FundTypeStock, RiskLevel: 3, Nav: 10000})
	f2, _ := svc.CreateFund(model.Fund{Code: "000011", Name: "F11", Type: model.FundTypeBond, RiskLevel: 2, Nav: 10000})
	acc, _ := svc.CreateAccount(model.Account{Owner: "J", Balance: 100000})
	p1, _ := svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: f1.ID, Amount: 30000, Interval: model.PlanIntervalMonthly, NextExecAt: time.Now().Add(-time.Hour)})
	p2, _ := svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: f2.ID, Amount: 10000, Interval: model.PlanIntervalMonthly, NextExecAt: time.Now().Add(-time.Hour)})
	svc.ExecutePlan(p1.ID)
	svc.ExecutePlan(p2.ID)

	top := svc.GetTopFunds(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 top funds, got %d", len(top))
	}
	if top[0].FundID != f1.ID {
		t.Fatalf("expected top fund %s, got %s", f1.ID, top[0].FundID)
	}
	if top[0].TotalInvested != 30000 {
		t.Fatalf("expected top invested 30000, got %d", top[0].TotalInvested)
	}
}

func TestFundTypeStats(t *testing.T) {
	svc := newTestService()
	svc.CreateFund(model.Fund{Code: "000012", Name: "F12", Type: model.FundTypeStock, RiskLevel: 3, Nav: 10000})
	svc.CreateFund(model.Fund{Code: "000013", Name: "F13", Type: model.FundTypeStock, RiskLevel: 3, Nav: 20000})
	svc.CreateFund(model.Fund{Code: "000014", Name: "F14", Type: model.FundTypeBond, RiskLevel: 2, Nav: 15000})

	stats := svc.GetFundTypeStats()
	if len(stats) != 2 {
		t.Fatalf("expected 2 type stats, got %d", len(stats))
	}
	for _, st := range stats {
		if st.Type == model.FundTypeStock && st.Count != 2 {
			t.Fatalf("expected stock count 2, got %d", st.Count)
		}
		if st.Type == model.FundTypeBond && st.Count != 1 {
			t.Fatalf("expected bond count 1, got %d", st.Count)
		}
	}
}

func TestFundRiskStats(t *testing.T) {
	svc := newTestService()
	svc.CreateFund(model.Fund{Code: "000015", Name: "F15", Type: model.FundTypeStock, RiskLevel: 3, Nav: 10000})
	svc.CreateFund(model.Fund{Code: "000016", Name: "F16", Type: model.FundTypeBond, RiskLevel: 2, Nav: 10000})

	stats := svc.GetFundRiskStats()
	if len(stats) != 2 {
		t.Fatalf("expected 2 risk stats, got %d", len(stats))
	}
}

func TestGetAccountHoldingSummary(t *testing.T) {
	svc := newTestService()
	fund, _ := svc.CreateFund(model.Fund{Code: "000017", Name: "F17", Type: model.FundTypeMixed, RiskLevel: 3, Nav: 10000})
	acc, _ := svc.CreateAccount(model.Account{Owner: "K", Balance: 50000})
	plan, _ := svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: fund.ID, Amount: 10000, Interval: model.PlanIntervalMonthly, NextExecAt: time.Now().Add(-time.Hour)})
	svc.ExecutePlan(plan.ID)

	summary, err := svc.GetAccountHoldingSummary(acc.ID)
	if err != nil {
		t.Fatalf("get account holding summary: %v", err)
	}
	if len(summary) != 1 {
		t.Fatalf("expected 1 summary item, got %d", len(summary))
	}
	if summary[0].FundID != fund.ID {
		t.Fatalf("unexpected fund id %s", summary[0].FundID)
	}
	if summary[0].Cost != 10000 {
		t.Fatalf("expected cost 10000, got %d", summary[0].Cost)
	}
}

func TestCreateDividendCrossEntity(t *testing.T) {
	svc := newTestService()
	_, err := svc.CreateDividend(model.Dividend{AccountID: "no", FundID: "no", Amount: 100, Type: model.DividendTypeCash})
	if err == nil {
		t.Fatal("expected error for missing account")
	}
	acc, _ := svc.CreateAccount(model.Account{Owner: "L", Balance: 10000})
	_, err = svc.CreateDividend(model.Dividend{AccountID: acc.ID, FundID: "no", Amount: 100, Type: model.DividendTypeCash})
	if err == nil {
		t.Fatal("expected error for missing fund")
	}
	fund, _ := svc.CreateFund(model.Fund{Code: "000018", Name: "F18", Type: model.FundTypeStock, RiskLevel: 3, Nav: 10000})
	div, err := svc.CreateDividend(model.Dividend{AccountID: acc.ID, FundID: fund.ID, Amount: 100, Type: model.DividendTypeReinvest})
	if err != nil {
		t.Fatalf("create dividend: %v", err)
	}
	if div.ID == "" {
		t.Fatal("dividend id empty")
	}
}

func TestListDividends(t *testing.T) {
	svc := newTestService()
	acc, _ := svc.CreateAccount(model.Account{Owner: "M", Balance: 10000})
	fund, _ := svc.CreateFund(model.Fund{Code: "000019", Name: "F19", Type: model.FundTypeBond, RiskLevel: 2, Nav: 10000})
	svc.CreateDividend(model.Dividend{AccountID: acc.ID, FundID: fund.ID, Amount: 100, Type: model.DividendTypeCash})
	items, total, err := svc.ListDividends(model.DividendFilter{AccountID: acc.ID}, 1, 10)
	if err != nil {
		t.Fatalf("list dividends: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
}

func TestExecutePlanNextExecAt(t *testing.T) {
	svc := newTestService()
	fund, _ := svc.CreateFund(model.Fund{Code: "000020", Name: "F20", Type: model.FundTypeStock, RiskLevel: 3, Nav: 10000})
	acc, _ := svc.CreateAccount(model.Account{Owner: "N", Balance: 100000})
	next := time.Now().Add(-time.Hour)
	plan, _ := svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: fund.ID, Amount: 10000, Interval: model.PlanIntervalDaily, NextExecAt: next})
	svc.ExecutePlan(plan.ID)
	plan2, _ := svc.GetPlan(plan.ID)
	if !plan2.NextExecAt.After(next) {
		t.Fatal("next_exec_at should be advanced")
	}
}

func TestPlanFilter(t *testing.T) {
	svc := newTestService()
	acc, _ := svc.CreateAccount(model.Account{Owner: "O", Balance: 10000})
	fund, _ := svc.CreateFund(model.Fund{Code: "000021", Name: "F21", Type: model.FundTypeIndex, RiskLevel: 2, Nav: 10000})
	p1, _ := svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: fund.ID, Amount: 1000, Interval: model.PlanIntervalMonthly})
	p2, _ := svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: fund.ID, Amount: 2000, Interval: model.PlanIntervalWeekly})

	items, total, _ := svc.ListPlans(model.PlanFilter{Interval: model.PlanIntervalMonthly}, 1, 10)
	if total != 1 || len(items) != 1 || items[0].ID != p1.ID {
		t.Fatalf("plan filter by interval failed")
	}
	items, total, _ = svc.ListPlans(model.PlanFilter{Interval: model.PlanIntervalWeekly}, 1, 10)
	if total != 1 || len(items) != 1 || items[0].ID != p2.ID {
		t.Fatalf("plan filter by interval failed for p2")
	}
	items, total, _ = svc.ListPlans(model.PlanFilter{Status: model.PlanStatusActive}, 1, 10)
	if total != 2 || len(items) != 2 {
		t.Fatalf("plan filter by status failed")
	}
	_, total, _ = svc.ListPlans(model.PlanFilter{AccountID: acc.ID}, 1, 10)
	if total != 2 {
		t.Fatalf("plan filter by account_id failed")
	}
}

func TestTransactionTimeFilter(t *testing.T) {
	svc := newTestService()
	fund, _ := svc.CreateFund(model.Fund{Code: "000022", Name: "F22", Type: model.FundTypeBond, RiskLevel: 2, Nav: 10000})
	acc, _ := svc.CreateAccount(model.Account{Owner: "P", Balance: 50000})
	plan, _ := svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: fund.ID, Amount: 5000, Interval: model.PlanIntervalDaily, NextExecAt: time.Now().Add(-time.Hour)})
	svc.ExecutePlan(plan.ID)

	start := time.Now().Add(-time.Minute)
	end := time.Now().Add(time.Minute)
	items, total, err := svc.ListTransactions(model.TransactionFilter{AccountID: acc.ID, StartTime: start, EndTime: end}, 1, 10)
	if err != nil {
		t.Fatalf("list transactions: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}
	items, total, _ = svc.ListTransactions(model.TransactionFilter{AccountID: acc.ID, StartTime: time.Now().Add(time.Hour)}, 1, 10)
	if total != 0 || len(items) != 0 {
		t.Fatalf("expected empty result for future start time")
	}
}

func TestFundFilterRiskLevel(t *testing.T) {
	svc := newTestService()
	svc.CreateFund(model.Fund{Code: "000023", Name: "F23", Type: model.FundTypeStock, RiskLevel: 1, Nav: 10000})
	svc.CreateFund(model.Fund{Code: "000024", Name: "F24", Type: model.FundTypeStock, RiskLevel: 3, Nav: 10000})
	items, total, _ := svc.ListFunds(model.FundFilter{RiskLevel: 1}, 1, 10)
	if total != 1 || len(items) != 1 {
		t.Fatalf("expected 1 fund with risk level 1, got %d", total)
	}
}

func TestUpdateFundConflict(t *testing.T) {
	svc := newTestService()
	f1, _ := svc.CreateFund(model.Fund{Code: "000025", Name: "F25", Type: model.FundTypeStock, RiskLevel: 3, Nav: 10000})
	f2, _ := svc.CreateFund(model.Fund{Code: "000026", Name: "F26", Type: model.FundTypeBond, RiskLevel: 2, Nav: 10000})
	_, err := svc.UpdateFund(f1.ID, model.Fund{Code: f2.Code})
	if err == nil {
		t.Fatal("expected conflict when updating fund code to existing code")
	}
}

func TestRedeemZeroHolding(t *testing.T) {
	svc := newTestService()
	fund, _ := svc.CreateFund(model.Fund{Code: "000027", Name: "F27", Type: model.FundTypeStock, RiskLevel: 3, Nav: 10000})
	acc, _ := svc.CreateAccount(model.Account{Owner: "Q", Balance: 100000})
	plan, _ := svc.CreatePlan(model.Plan{AccountID: acc.ID, FundID: fund.ID, Amount: 10000, Interval: model.PlanIntervalMonthly, NextExecAt: time.Now().Add(-time.Hour)})
	svc.ExecutePlan(plan.ID)
	shares := int64((10000 * 10000) / 10000)
	svc.Redeem(acc.ID, fund.ID, shares)
	holdings, _, _ := svc.ListHoldings(model.HoldingFilter{AccountID: acc.ID}, 1, 10)
	if len(holdings) != 1 {
		t.Fatalf("expected 1 holding")
	}
	if holdings[0].TotalShares != 0 {
		t.Fatalf("expected 0 shares after full redeem, got %d", holdings[0].TotalShares)
	}
	if holdings[0].Cost != 0 {
		t.Fatalf("expected 0 cost after full redeem, got %d", holdings[0].Cost)
	}
}

func TestGetTopFundsEmpty(t *testing.T) {
	svc := newTestService()
	top := svc.GetTopFunds(5)
	if len(top) != 0 {
		t.Fatalf("expected empty top funds, got %d", len(top))
	}
}

func TestPagination(t *testing.T) {
	svc := newTestService()
	for i := 0; i < 25; i++ {
		code := "C" + strconv.Itoa(i)
		svc.CreateFund(model.Fund{Code: code, Name: "F", Type: model.FundTypeStock, RiskLevel: 3, Nav: 10000})
	}
	items, total, _ := svc.ListFunds(model.FundFilter{}, 1, 10)
	if total != 25 || len(items) != 10 {
		t.Fatalf("page 1 expected 10 items, got %d", len(items))
	}
	items, total, _ = svc.ListFunds(model.FundFilter{}, 3, 10)
	if total != 25 || len(items) != 5 {
		t.Fatalf("page 3 expected 5 items, got %d", len(items))
	}
}

