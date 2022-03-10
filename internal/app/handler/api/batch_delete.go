package api

import (
	"errors"
	"net/http"
	"shortener/internal/app/handler"
	"shortener/internal/app/service/store"
)

// BatchRemoveHandler removes multiple urls.
//
//	curl -v -X DELETE -H "Content-Type: application/json" -d '["xxx", "xxy"]' http://localhost:8080/api/user/urls
func BatchRemoveHandler(s store.BatchRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := make([]string, 0)
		err := readBody(r, &req)
		if err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}

		uid := handler.ReadContextString(r.Context(), handler.ContextKeyUID{})

		if err := s.BatchRemove(uid, req...); err != nil {
			if errors.Is(err, store.ErrBadInput) {
				writeError(w, err, http.StatusBadRequest)
			} else {
				writeError(w, err, http.StatusInternalServerError)
			}
			return
		}

		writeResponse(w, nil, http.StatusAccepted)
	}
}
