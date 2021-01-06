package server

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) HandleMonthRequest(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	monthParam := params.ByName("month")
	if monthParam == "" {
		s.error(w, r, "month path param missing", 400)
	}

	monthNumber, err := strconv.Atoi(monthParam)
	if err != nil {
		s.error(w, r, "month is not a number", 400)
	}

	if monthNumber < 1 || monthNumber > 12 {
		s.error(w, r, "month must be in the range 1-12", 400)
	}

	yearParam := params.ByName("year")
	if monthParam == "" {
		s.error(w, r, "year path param missing", 400)
	}

	yearNumber, err := strconv.Atoi(yearParam)
	if err != nil {
		s.error(w, r, "year is not a number", 400)
	}

	month := time.Date(yearNumber, time.Month(monthNumber), 1, 0, 0, 0, 0, time.UTC)

	result, err := s.storage.GetMonths(month, month)
	if err != nil {
		s.error(w, r, "internal server error", 500)
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Transfer-Encoding", "identity")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		s.error(w, r, "internal server error", 500)
	}
}
