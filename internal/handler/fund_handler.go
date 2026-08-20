package handler

import (
	"net/http"

	"fundinvest/internal/model"
	"fundinvest/pkg/httpx"
)

func (s *Server) registerFundRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/funds", s.createFund)
	mux.HandleFunc("GET /api/funds", s.listFunds)
	mux.HandleFunc("GET /api/funds/{id}", s.getFund)
	mux.HandleFunc("PUT /api/funds/{id}", s.updateFund)
	mux.HandleFunc("DELETE /api/funds/{id}", s.deleteFund)
}

type createFundRequest struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	RiskLevel int    `json:"risk_level"`
	Nav       int64  `json:"nav"`
}

func (s *Server) createFund(w http.ResponseWriter, r *http.Request) {
	var req createFundRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.CreateFund(model.Fund{Code: req.Code, Name: req.Name, Type: req.Type, RiskLevel: req.RiskLevel, Nav: req.Nav})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, f)
}

func (s *Server) listFunds(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.FundFilter{
		Type:      r.URL.Query().Get("type"),
		Status:    r.URL.Query().Get("status"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListFunds(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getFund(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	f, err := s.svc.GetFund(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

type updateFundRequest struct {
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	RiskLevel int    `json:"risk_level,omitempty"`
	Nav       int64  `json:"nav,omitempty"`
	Status    string `json:"status,omitempty"`
}

func (s *Server) updateFund(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateFundRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.UpdateFund(id, model.Fund{Name: req.Name, Type: req.Type, RiskLevel: req.RiskLevel, Nav: req.Nav, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

func (s *Server) deleteFund(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteFund(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
