package handler

import (
	"net/http"

	"fundinvest/internal/model"
	"fundinvest/pkg/httpx"
)

func (s *Server) registerDividendRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/dividends", s.createDividend)
	mux.HandleFunc("GET /api/dividends", s.listDividends)
	mux.HandleFunc("GET /api/dividends/{id}", s.getDividend)
	mux.HandleFunc("PUT /api/dividends/{id}", s.updateDividend)
	mux.HandleFunc("DELETE /api/dividends/{id}", s.deleteDividend)
}

type createDividendRequest struct {
	AccountID string `json:"account_id"`
	FundID    string `json:"fund_id"`
	Amount    int64  `json:"amount"`
	Type      string `json:"type"`
}

func (s *Server) createDividend(w http.ResponseWriter, r *http.Request) {
	var req createDividendRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.CreateDividend(model.Dividend{AccountID: req.AccountID, FundID: req.FundID, Amount: req.Amount, Type: req.Type})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, d)
}

func (s *Server) listDividends(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.DividendFilter{
		AccountID: r.URL.Query().Get("account_id"),
		FundID:    r.URL.Query().Get("fund_id"),
		Type:      r.URL.Query().Get("type"),
	}
	items, total, err := s.svc.ListDividends(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getDividend(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	d, err := s.svc.GetDividend(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

type updateDividendRequest struct {
	Amount int64  `json:"amount,omitempty"`
	Type   string `json:"type,omitempty"`
}

func (s *Server) updateDividend(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateDividendRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.UpdateDividend(id, model.Dividend{Amount: req.Amount, Type: req.Type})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) deleteDividend(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteDividend(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
