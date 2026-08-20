package handler

import (
	"net/http"

	"fundinvest/internal/model"
	"fundinvest/pkg/httpx"
)

func (s *Server) registerAccountRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/accounts", s.createAccount)
	mux.HandleFunc("GET /api/accounts", s.listAccounts)
	mux.HandleFunc("GET /api/accounts/{id}", s.getAccount)
	mux.HandleFunc("PUT /api/accounts/{id}", s.updateAccount)
	mux.HandleFunc("DELETE /api/accounts/{id}", s.deleteAccount)
}

type createAccountRequest struct {
	Owner   string `json:"owner"`
	Balance int64  `json:"balance"`
}

func (s *Server) createAccount(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateAccount(model.Account{Owner: req.Owner, Balance: req.Balance})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listAccounts(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AccountFilter{Keyword: r.URL.Query().Get("keyword")}
	items, total, err := s.svc.ListAccounts(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAccount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.GetAccount(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

type updateAccountRequest struct {
	Owner   string `json:"owner,omitempty"`
	Balance int64  `json:"balance,omitempty"`
}

func (s *Server) updateAccount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateAccountRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.UpdateAccount(id, model.Account{Owner: req.Owner, Balance: req.Balance})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteAccount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAccount(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
