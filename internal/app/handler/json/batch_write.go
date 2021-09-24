package json

import (
	"errors"
	"github.com/russianlagman/go-musthave-shortener/internal/app/handler"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store"
	"net/http"
)

type BatchWriteRequestItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchWriteResponseItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func BatchWriteHandler(s store.BatchWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqObj := make([]BatchWriteRequestItem, 0)
		err := readBody(r, &reqObj)
		if err != nil {
			writeError(w, err, http.StatusBadRequest)
			return
		}

		uid := handler.ReadContextString(r.Context(), handler.ContextKeyUID{})

		storeReq := make([]store.Record, len(reqObj))
		for i, rec := range reqObj {
			storeReq[i] = store.Record{
				CorrelationID: rec.CorrelationID,
				OriginalURL:   rec.OriginalURL,
			}
		}

		storeRes, err := s.BatchWrite(uid, storeReq)
		if err != nil {
			if errors.Is(err, store.ErrBadInput) {
				writeError(w, err, http.StatusBadRequest)
			} else {
				writeError(w, err, http.StatusInternalServerError)
			}
			return
		}

		respObj := make([]BatchWriteResponseItem, len(reqObj))
		for i, rec := range storeRes {
			respObj[i] = BatchWriteResponseItem{
				CorrelationID: rec.CorrelationID,
				ShortURL:      rec.ShortURL,
			}
		}

		writeResponse(w, respObj, http.StatusCreated)
	}
}
