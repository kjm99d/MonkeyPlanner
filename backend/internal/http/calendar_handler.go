package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/kjm99d/monkey-planner/backend/internal/domain"
	"github.com/kjm99d/monkey-planner/backend/internal/service"
	"github.com/kjm99d/monkey-planner/backend/internal/storage"
)

type calendarHandler struct{ svc *service.Service }

// GET /api/calendar?year=YYYY&month=MM
func (h *calendarHandler) month(w http.ResponseWriter, r *http.Request) {
	year, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil || year < 1970 || year > 9999 {
		writeErr(w, http.StatusBadRequest, "invalid_year", "year must be 1970–9999")
		return
	}
	monthNum, err := strconv.Atoi(r.URL.Query().Get("month"))
	if err != nil || monthNum < 1 || monthNum > 12 {
		writeErr(w, http.StatusBadRequest, "invalid_month", "month must be 1–12")
		return
	}
	out, err := h.svc.GetMonthStats(r.Context(), year, time.Month(monthNum))
	if err != nil {
		mapError(w, err)
		return
	}
	if out == nil {
		out = []storage.DayCount{}
	}
	writeJSON(w, http.StatusOK, out)
}

// GET /api/calendar/day?date=YYYY-MM-DD
func (h *calendarHandler) day(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	day, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_date", "date must be YYYY-MM-DD")
		return
	}
	out, err := h.svc.GetDayStats(r.Context(), day.UTC())
	if err != nil {
		mapError(w, err)
		return
	}
	// nil slice → [] 로 정규화 (클라이언트가 .length 접근 시 TypeError 방지)
	if out.Created == nil {
		out.Created = []domain.Issue{}
	}
	if out.Approved == nil {
		out.Approved = []domain.Issue{}
	}
	if out.Completed == nil {
		out.Completed = []domain.Issue{}
	}
	writeJSON(w, http.StatusOK, out)
}
