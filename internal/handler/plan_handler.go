package handler

import (
	"net/http"
	"time"

	"fundinvest/internal/model"
	"fundinvest/pkg/httpx"
)

func (s *Server) registerPlanRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/plans", s.createPlan)
	mux.HandleFunc("GET /api/plans", s.listPlans)
	mux.HandleFunc("GET /api/plans/{id}", s.getPlan)
	mux.HandleFunc("PUT /api/plans/{id}", s.updatePlan)
	mux.HandleFunc("POST /api/plans/{id}/status", s.changePlanStatus)
	mux.HandleFunc("POST /api/plans/{id}/execute", s.executePlan)
	mux.HandleFunc("DELETE /api/plans/{id}", s.deletePlan)
}

type createPlanRequest struct {
	AccountID  string    `json:"account_id"`
	FundID     string    `json:"fund_id"`
	Amount     int64     `json:"amount"`
	Interval   string    `json:"interval"`
	NextExecAt time.Time `json:"next_exec_at"`
}

func (s *Server) createPlan(w http.ResponseWriter, r *http.Request) {
	var req createPlanRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.CreatePlan(model.Plan{AccountID: req.AccountID, FundID: req.FundID, Amount: req.Amount, Interval: req.Interval, NextExecAt: req.NextExecAt})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, p)
}

func (s *Server) listPlans(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.PlanFilter{
		AccountID: r.URL.Query().Get("account_id"),
		FundID:    r.URL.Query().Get("fund_id"),
		Status:    r.URL.Query().Get("status"),
		Interval:  r.URL.Query().Get("interval"),
	}
	items, total, err := s.svc.ListPlans(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getPlan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := s.svc.GetPlan(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

type updatePlanRequest struct {
	Amount     int64     `json:"amount,omitempty"`
	Interval   string    `json:"interval,omitempty"`
	NextExecAt time.Time `json:"next_exec_at,omitempty"`
}

func (s *Server) updatePlan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updatePlanRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.UpdatePlan(id, model.Plan{Amount: req.Amount, Interval: req.Interval, NextExecAt: req.NextExecAt})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

type changePlanStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) changePlanStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req changePlanStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.ChangePlanStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) executePlan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	txn, err := s.svc.ExecutePlan(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, txn)
}

func (s *Server) deletePlan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeletePlan(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
