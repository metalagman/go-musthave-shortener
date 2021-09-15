package json

import (
	"errors"
	"github.com/russianlagman/go-musthave-shortener/internal/app/middleware"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store"
	"net/http"
)

type WriteHandlerRequest struct {
	URL string `json:"url"`
}

type WriteHandlerResponse struct {
	Result string `json:"result"`
}

func WriteHandler(s store.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqObj := &WriteHandlerRequest{}
		err := readBody(r, reqObj)
		if err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}

		shortURL, err := s.WriteUserURL(reqObj.URL, r.Context().Value(middleware.ContextKeyUID{}).(string))
		if err != nil {
			if errors.Is(err, store.ErrBadInput) {
				writeError(w, err, http.StatusBadRequest)
			} else {
				writeError(w, err, http.StatusInternalServerError)
			}
			return
		}

		respObj := &WriteHandlerResponse{Result: shortURL}
		writeResponse(w, respObj, http.StatusCreated)
	}
}
