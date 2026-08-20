package handler

import (
	"net/http"

	"fundinvest/internal/model"
	"fundinvest/pkg/httpx"
)

func (s *Server) registerTransactionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/transactions", s.createTransaction)
	mux.HandleFunc("GET /api/transactions", s.listTransactions)
	mux.HandleFunc("GET /api/transactions/{id}", s.getTransaction)
	mux.HandleFunc("PUT /api/transactions/{id}", s.updateTransaction)
	mux.HandleFunc("POST /api/transactions/redeem", s.redeem)
	mux.HandleFunc("DELETE /api/transactions/{id}", s.deleteTransaction)
}

type createTransactionRequest struct {
	AccountID string `json:"account_id"`
	FundID    string `json:"fund_id"`
	PlanID    string `json:"plan_id,omitempty"`
	Type      string `json:"type"`
	Amount    int64  `json:"amount"`
	Shares    int64  `json:"shares"`
	Price     int64  `json:"price"`
}

func (s *Server) createTransaction(w http.ResponseWriter, r *http.Request) {
	var req createTransactionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateTransaction(model.Transaction{AccountID: req.AccountID, FundID: req.FundID, PlanID: req.PlanID, Type: req.Type, Amount: req.Amount, Shares: req.Shares, Price: req.Price})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listTransactions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TransactionFilter{
		AccountID: r.URL.Query().Get("account_id"),
		FundID:    r.URL.Query().Get("fund_id"),
		Type:      r.URL.Query().Get("type"),
		Status:    r.URL.Query().Get("status"),
		PlanID:    r.URL.Query().Get("plan_id"),
	}
	items, total, err := s.svc.ListTransactions(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.GetTransaction(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type updateTransactionRequest struct {
	Status string `json:"status,omitempty"`
	Amount int64  `json:"amount,omitempty"`
	Shares int64  `json:"shares,omitempty"`
	Price  int64  `json:"price,omitempty"`
}

func (s *Server) updateTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateTransactionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateTransaction(id, model.Transaction{Status: req.Status, Amount: req.Amount, Shares: req.Shares, Price: req.Price})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type redeemRequest struct {
	AccountID string `json:"account_id"`
	FundID    string `json:"fund_id"`
	Shares    int64  `json:"shares"`
}

func (s *Server) redeem(w http.ResponseWriter, r *http.Request) {
	var req redeemRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	txn, err := s.svc.Redeem(req.AccountID, req.FundID, req.Shares)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, txn)
}

func (s *Server) deleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteTransaction(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
