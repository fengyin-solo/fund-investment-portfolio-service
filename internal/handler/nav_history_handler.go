package handler

import (
	"net/http"
	"time"

	"fundinvest/internal/model"
	"fundinvest/pkg/httpx"
)

func (s *Server) registerNavHistoryRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/nav-history", s.createNavHistory)
	mux.HandleFunc("GET /api/nav-history", s.listNavHistory)
	mux.HandleFunc("GET /api/nav-history/{id}", s.getNavHistory)
	mux.HandleFunc("PUT /api/nav-history/{id}", s.updateNavHistory)
	mux.HandleFunc("DELETE /api/nav-history/{id}", s.deleteNavHistory)
}

type createNavHistoryRequest struct {
	FundID     string    `json:"fund_id"`
	Nav        int64     `json:"nav"`
	RecordedAt time.Time `json:"recorded_at,omitempty"`
}

func (s *Server) createNavHistory(w http.ResponseWriter, r *http.Request) {
	var req createNavHistoryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.CreateNavHistory(model.NavHistory{FundID: req.FundID, Nav: req.Nav, RecordedAt: req.RecordedAt})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, n)
}

func (s *Server) listNavHistory(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.NavHistoryFilter{
		FundID: r.URL.Query().Get("fund_id"),
	}
	items, total, err := s.svc.ListNavHistory(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getNavHistory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, err := s.svc.GetNavHistory(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

type updateNavHistoryRequest struct {
	Nav        int64     `json:"nav,omitempty"`
	RecordedAt time.Time `json:"recorded_at,omitempty"`
}

func (s *Server) updateNavHistory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateNavHistoryRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.UpdateNavHistory(id, model.NavHistory{Nav: req.Nav, RecordedAt: req.RecordedAt})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

func (s *Server) deleteNavHistory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteNavHistory(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
