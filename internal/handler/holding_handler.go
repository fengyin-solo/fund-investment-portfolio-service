package handler

import (
	"net/http"

	"fundinvest/internal/model"
	"fundinvest/pkg/httpx"
)

func (s *Server) registerHoldingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/holdings", s.createHolding)
	mux.HandleFunc("GET /api/holdings", s.listHoldings)
	mux.HandleFunc("GET /api/holdings/{id}", s.getHolding)
	mux.HandleFunc("PUT /api/holdings/{id}", s.updateHolding)
	mux.HandleFunc("DELETE /api/holdings/{id}", s.deleteHolding)
}

type createHoldingRequest struct {
	AccountID   string `json:"account_id"`
	FundID      string `json:"fund_id"`
	TotalShares int64  `json:"total_shares"`
	Cost        int64  `json:"cost"`
	MarketValue int64  `json:"market_value"`
}

func (s *Server) createHolding(w http.ResponseWriter, r *http.Request) {
	var req createHoldingRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	h, err := s.svc.CreateHolding(model.Holding{AccountID: req.AccountID, FundID: req.FundID, TotalShares: req.TotalShares, Cost: req.Cost, MarketValue: req.MarketValue})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, h)
}

func (s *Server) listHoldings(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.HoldingFilter{
		AccountID: r.URL.Query().Get("account_id"),
		FundID:    r.URL.Query().Get("fund_id"),
	}
	items, total, err := s.svc.ListHoldings(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getHolding(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	h, err := s.svc.GetHolding(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, h)
}

type updateHoldingRequest struct {
	TotalShares int64 `json:"total_shares,omitempty"`
	Cost        int64 `json:"cost,omitempty"`
	MarketValue int64 `json:"market_value,omitempty"`
}

func (s *Server) updateHolding(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateHoldingRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	h, err := s.svc.UpdateHolding(id, model.Holding{TotalShares: req.TotalShares, Cost: req.Cost, MarketValue: req.MarketValue})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, h)
}

func (s *Server) deleteHolding(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteHolding(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
