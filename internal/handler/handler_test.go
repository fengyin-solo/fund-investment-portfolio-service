package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"fundinvest/internal/config"
	"fundinvest/internal/service"
	"fundinvest/internal/store"
	"fundinvest/pkg/httpx"
	"fundinvest/pkg/logger"
)

func newTestServer() *httptest.Server {
	cfg := &config.Config{Addr: ":8080", MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	st := store.NewMemoryStore()
	svc := service.New(st, log, cfg)
	server := NewServer(svc, log, cfg)
	return httptest.NewServer(server.Routes())
}

func TestCreateFundHandler(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	body, _ := json.Marshal(map[string]interface{}{"code": "H001", "name": "HF", "type": "stock", "risk_level": 3, "nav": 15000})
	resp, err := http.Post(srv.URL+"/api/funds", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("post fund: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var r httpx.Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if r.Code != 0 {
		t.Fatalf("expected code 0, got %d", r.Code)
	}
}

func TestCreateFundValidationError(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	body, _ := json.Marshal(map[string]interface{}{"code": "", "name": "", "type": "", "risk_level": 0, "nav": 0})
	resp, err := http.Post(srv.URL+"/api/funds", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("post fund: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestGetFundNotFound(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/api/funds/notexist")
	if err != nil {
		t.Fatalf("get fund: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestCreateAccountHandler(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	body, _ := json.Marshal(map[string]interface{}{"owner": "Test", "balance": 50000})
	resp, err := http.Post(srv.URL+"/api/accounts", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("post account: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestListAccountsPagination(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	for i := 0; i < 5; i++ {
		body, _ := json.Marshal(map[string]interface{}{"owner": strings.Repeat("A", i+1), "balance": 1000})
		http.Post(srv.URL+"/api/accounts", "application/json", bytes.NewReader(body))
	}
	resp, err := http.Get(srv.URL + "/api/accounts?page=1&size=2")
	if err != nil {
		t.Fatalf("list accounts: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var r httpx.Response
	json.NewDecoder(resp.Body).Decode(&r)
	m, ok := r.Data.(map[string]interface{})
	if !ok {
		t.Fatal("expected map data")
	}
	items, ok := m["items"].([]interface{})
	if !ok || len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
}

func TestCreatePlanHandler(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	body, _ := json.Marshal(map[string]interface{}{"code": "H002", "name": "HF2", "type": "bond", "risk_level": 2, "nav": 10000})
	resp, _ := http.Post(srv.URL+"/api/funds", "application/json", bytes.NewReader(body))
	var fundR httpx.Response
	json.NewDecoder(resp.Body).Decode(&fundR)
	resp.Body.Close()
	fund := fundR.Data.(map[string]interface{})
	fundID := fund["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"owner": "P", "balance": 100000})
	resp, _ = http.Post(srv.URL+"/api/accounts", "application/json", bytes.NewReader(body))
	var accR httpx.Response
	json.NewDecoder(resp.Body).Decode(&accR)
	resp.Body.Close()
	acc := accR.Data.(map[string]interface{})
	accID := acc["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"account_id": accID, "fund_id": fundID, "amount": 5000, "interval": "monthly", "next_exec_at": time.Now().Add(-time.Hour).Format(time.RFC3339)})
	resp, err := http.Post(srv.URL+"/api/plans", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("post plan: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestExecutePlanHandler(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	body, _ := json.Marshal(map[string]interface{}{"code": "H003", "name": "HF3", "type": "stock", "risk_level": 3, "nav": 10000})
	resp, _ := http.Post(srv.URL+"/api/funds", "application/json", bytes.NewReader(body))
	var fundR httpx.Response
	json.NewDecoder(resp.Body).Decode(&fundR)
	resp.Body.Close()
	fundID := fundR.Data.(map[string]interface{})["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"owner": "E", "balance": 100000})
	resp, _ = http.Post(srv.URL+"/api/accounts", "application/json", bytes.NewReader(body))
	var accR httpx.Response
	json.NewDecoder(resp.Body).Decode(&accR)
	resp.Body.Close()
	accID := accR.Data.(map[string]interface{})["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"account_id": accID, "fund_id": fundID, "amount": 5000, "interval": "monthly", "next_exec_at": time.Now().Add(-time.Hour).Format(time.RFC3339)})
	resp, _ = http.Post(srv.URL+"/api/plans", "application/json", bytes.NewReader(body))
	var planR httpx.Response
	json.NewDecoder(resp.Body).Decode(&planR)
	resp.Body.Close()
	planID := planR.Data.(map[string]interface{})["id"].(string)

	resp, err := http.Post(srv.URL+"/api/plans/"+planID+"/execute", "application/json", nil)
	if err != nil {
		t.Fatalf("execute plan: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRedeemHandler(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	body, _ := json.Marshal(map[string]interface{}{"code": "H004", "name": "HF4", "type": "mixed", "risk_level": 3, "nav": 10000})
	resp, _ := http.Post(srv.URL+"/api/funds", "application/json", bytes.NewReader(body))
	var fundR httpx.Response
	json.NewDecoder(resp.Body).Decode(&fundR)
	resp.Body.Close()
	fundID := fundR.Data.(map[string]interface{})["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"owner": "R", "balance": 100000})
	resp, _ = http.Post(srv.URL+"/api/accounts", "application/json", bytes.NewReader(body))
	var accR httpx.Response
	json.NewDecoder(resp.Body).Decode(&accR)
	resp.Body.Close()
	accID := accR.Data.(map[string]interface{})["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"account_id": accID, "fund_id": fundID, "amount": 20000, "interval": "weekly", "next_exec_at": time.Now().Add(-time.Hour).Format(time.RFC3339)})
	resp, _ = http.Post(srv.URL+"/api/plans", "application/json", bytes.NewReader(body))
	var planR httpx.Response
	json.NewDecoder(resp.Body).Decode(&planR)
	resp.Body.Close()
	planID := planR.Data.(map[string]interface{})["id"].(string)

	http.Post(srv.URL+"/api/plans/"+planID+"/execute", "application/json", nil)

	body, _ = json.Marshal(map[string]interface{}{"account_id": accID, "fund_id": fundID, "shares": 5000})
	resp, err := http.Post(srv.URL+"/api/transactions/redeem", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("redeem: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestStatsHandlers(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	body, _ := json.Marshal(map[string]interface{}{"code": "H005", "name": "HF5", "type": "stock", "risk_level": 3, "nav": 10000})
	resp, _ := http.Post(srv.URL+"/api/funds", "application/json", bytes.NewReader(body))
	var fundR httpx.Response
	json.NewDecoder(resp.Body).Decode(&fundR)
	resp.Body.Close()
	fundID := fundR.Data.(map[string]interface{})["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"owner": "S", "balance": 100000})
	resp, _ = http.Post(srv.URL+"/api/accounts", "application/json", bytes.NewReader(body))
	var accR httpx.Response
	json.NewDecoder(resp.Body).Decode(&accR)
	resp.Body.Close()
	accID := accR.Data.(map[string]interface{})["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"account_id": accID, "fund_id": fundID, "amount": 10000, "interval": "monthly", "next_exec_at": time.Now().Add(-time.Hour).Format(time.RFC3339)})
	resp, _ = http.Post(srv.URL+"/api/plans", "application/json", bytes.NewReader(body))
	var planR httpx.Response
	json.NewDecoder(resp.Body).Decode(&planR)
	resp.Body.Close()
	planID := planR.Data.(map[string]interface{})["id"].(string)

	http.Post(srv.URL+"/api/plans/"+planID+"/execute", "application/json", nil)

	resp, err := http.Get(srv.URL + "/api/stats/holdings/" + accID)
	if err != nil {
		t.Fatalf("get holdings stats: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	resp, err = http.Get(srv.URL + "/api/stats/profit/" + accID)
	if err != nil {
		t.Fatalf("get profit stats: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	resp, err = http.Get(srv.URL + "/api/stats/fund-types")
	if err != nil {
		t.Fatalf("get fund type stats: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	resp, err = http.Get(srv.URL + "/api/stats/fund-risks")
	if err != nil {
		t.Fatalf("get fund risk stats: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	resp, err = http.Get(srv.URL + "/api/stats/top-funds?n=5")
	if err != nil {
		t.Fatalf("get top funds: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestCreateTransactionHandler(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	body, _ := json.Marshal(map[string]interface{}{"code": "H006", "name": "HF6", "type": "bond", "risk_level": 2, "nav": 10000})
	resp, _ := http.Post(srv.URL+"/api/funds", "application/json", bytes.NewReader(body))
	var fundR httpx.Response
	json.NewDecoder(resp.Body).Decode(&fundR)
	resp.Body.Close()
	fundID := fundR.Data.(map[string]interface{})["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"owner": "T", "balance": 50000})
	resp, _ = http.Post(srv.URL+"/api/accounts", "application/json", bytes.NewReader(body))
	var accR httpx.Response
	json.NewDecoder(resp.Body).Decode(&accR)
	resp.Body.Close()
	accID := accR.Data.(map[string]interface{})["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"account_id": accID, "fund_id": fundID, "type": "invest", "amount": 5000, "shares": 5000, "price": 10000})
	resp, err := http.Post(srv.URL+"/api/transactions", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("post transaction: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestPlanStatusHandler(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	body, _ := json.Marshal(map[string]interface{}{"code": "H007", "name": "HF7", "type": "index", "risk_level": 2, "nav": 10000})
	resp, _ := http.Post(srv.URL+"/api/funds", "application/json", bytes.NewReader(body))
	var fundR httpx.Response
	json.NewDecoder(resp.Body).Decode(&fundR)
	resp.Body.Close()
	fundID := fundR.Data.(map[string]interface{})["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"owner": "U", "balance": 10000})
	resp, _ = http.Post(srv.URL+"/api/accounts", "application/json", bytes.NewReader(body))
	var accR httpx.Response
	json.NewDecoder(resp.Body).Decode(&accR)
	resp.Body.Close()
	accID := accR.Data.(map[string]interface{})["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"account_id": accID, "fund_id": fundID, "amount": 1000, "interval": "monthly"})
	resp, _ = http.Post(srv.URL+"/api/plans", "application/json", bytes.NewReader(body))
	var planR httpx.Response
	json.NewDecoder(resp.Body).Decode(&planR)
	resp.Body.Close()
	planID := planR.Data.(map[string]interface{})["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"status": "paused"})
	resp, err := http.Post(srv.URL+"/api/plans/"+planID+"/status", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("change status: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestDeleteFundHandler(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	body, _ := json.Marshal(map[string]interface{}{"code": "H008", "name": "HF8", "type": "stock", "risk_level": 3, "nav": 10000})
	resp, _ := http.Post(srv.URL+"/api/funds", "application/json", bytes.NewReader(body))
	var fundR httpx.Response
	json.NewDecoder(resp.Body).Decode(&fundR)
	resp.Body.Close()
	fundID := fundR.Data.(map[string]interface{})["id"].(string)

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/api/funds/"+fundID, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("delete fund: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}
}

func TestCreateHoldingHandler(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	body, _ := json.Marshal(map[string]interface{}{"code": "H009", "name": "HF9", "type": "bond", "risk_level": 2, "nav": 10000})
	resp, _ := http.Post(srv.URL+"/api/funds", "application/json", bytes.NewReader(body))
	var fundR httpx.Response
	json.NewDecoder(resp.Body).Decode(&fundR)
	resp.Body.Close()
	fundID := fundR.Data.(map[string]interface{})["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"owner": "V", "balance": 10000})
	resp, _ = http.Post(srv.URL+"/api/accounts", "application/json", bytes.NewReader(body))
	var accR httpx.Response
	json.NewDecoder(resp.Body).Decode(&accR)
	resp.Body.Close()
	accID := accR.Data.(map[string]interface{})["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"account_id": accID, "fund_id": fundID, "total_shares": 1000, "cost": 5000, "market_value": 6000})
	resp, err := http.Post(srv.URL+"/api/holdings", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("post holding: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestCreateNavHistoryHandler(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	body, _ := json.Marshal(map[string]interface{}{"code": "H010", "name": "HF10", "type": "stock", "risk_level": 3, "nav": 10000})
	resp, _ := http.Post(srv.URL+"/api/funds", "application/json", bytes.NewReader(body))
	var fundR httpx.Response
	json.NewDecoder(resp.Body).Decode(&fundR)
	resp.Body.Close()
	fundID := fundR.Data.(map[string]interface{})["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"fund_id": fundID, "nav": 10500})
	resp, err := http.Post(srv.URL+"/api/nav-history", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("post nav history: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestCreateDividendHandler(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	body, _ := json.Marshal(map[string]interface{}{"code": "H011", "name": "HF11", "type": "mixed", "risk_level": 3, "nav": 10000})
	resp, _ := http.Post(srv.URL+"/api/funds", "application/json", bytes.NewReader(body))
	var fundR httpx.Response
	json.NewDecoder(resp.Body).Decode(&fundR)
	resp.Body.Close()
	fundID := fundR.Data.(map[string]interface{})["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"owner": "W", "balance": 10000})
	resp, _ = http.Post(srv.URL+"/api/accounts", "application/json", bytes.NewReader(body))
	var accR httpx.Response
	json.NewDecoder(resp.Body).Decode(&accR)
	resp.Body.Close()
	accID := accR.Data.(map[string]interface{})["id"].(string)

	body, _ = json.Marshal(map[string]interface{}{"account_id": accID, "fund_id": fundID, "amount": 500, "type": "cash"})
	resp, err := http.Post(srv.URL+"/api/dividends", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("post dividend: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestInvalidJSONBody(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	resp, err := http.Post(srv.URL+"/api/funds", "application/json", bytes.NewReader([]byte("not json")))
	if err != nil {
		t.Fatalf("post invalid: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestPlanCrossEntityValidationHandler(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	body, _ := json.Marshal(map[string]interface{}{"account_id": "notexist", "fund_id": "notexist", "amount": 1000, "interval": "monthly"})
	resp, err := http.Post(srv.URL+"/api/plans", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("post plan: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing account, got %d", resp.StatusCode)
	}
}

func TestTransactionFilterHandler(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/api/transactions?type=invest&status=confirmed")
	if err != nil {
		t.Fatalf("get transactions: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
