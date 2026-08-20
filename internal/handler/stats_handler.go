package handler

import (
	"net/http"
	"strconv"

	"fundinvest/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/holdings/{account_id}", s.getAccountHoldingSummary)
	mux.HandleFunc("GET /api/stats/profit/{account_id}", s.getAccountProfit)
	mux.HandleFunc("GET /api/stats/fund-types", s.getFundTypeStats)
	mux.HandleFunc("GET /api/stats/fund-risks", s.getFundRiskStats)
	mux.HandleFunc("GET /api/stats/top-funds", s.getTopFunds)
}

func (s *Server) getAccountHoldingSummary(w http.ResponseWriter, r *http.Request) {
	accountID := r.PathValue("account_id")
	result, err := s.svc.GetAccountHoldingSummary(accountID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) getAccountProfit(w http.ResponseWriter, r *http.Request) {
	accountID := r.PathValue("account_id")
	result, err := s.svc.GetAccountProfit(accountID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) getFundTypeStats(w http.ResponseWriter, r *http.Request) {
	result := s.svc.GetFundTypeStats()
	httpx.OK(w, result)
}

func (s *Server) getFundRiskStats(w http.ResponseWriter, r *http.Request) {
	result := s.svc.GetFundRiskStats()
	httpx.OK(w, result)
}

func (s *Server) getTopFunds(w http.ResponseWriter, r *http.Request) {
	n, _ := strconv.Atoi(r.URL.Query().Get("n"))
	if n <= 0 {
		n = 10
	}
	result := s.svc.GetTopFunds(n)
	httpx.OK(w, result)
}
