package api

import (
	"net/http"
	"shortener/internal/app/service/store"
)

type StatResponse struct {
	URLCount  int `json:"urls"`
	UserCount int `json:"users"`
}

// StatHandler displays statistics
//
//	curl -X GET -H "Content-Type: application/json" --cookie "uid=XXX" http://localhost:8080/api/internal/stats
func StatHandler(s store.StatProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		res, err := s.Stat()
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}

		respObj := &StatResponse{
			UserCount: res.UserCount,
			URLCount:  res.UserCount,
		}

		writeResponse(w, respObj, http.StatusOK)
	}
}
