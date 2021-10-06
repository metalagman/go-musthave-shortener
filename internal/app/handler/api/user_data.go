package api

import (
	"net/http"
	"shortener/internal/app/handler"
	"shortener/internal/app/service/store"
)

type UserDataResponse []UserDataItem

type UserDataItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func UserDataHandler(s store.UserDataReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		rows := s.ReadUserData(handler.ReadContextString(r.Context(), handler.ContextKeyUID{}))
		respObj := make(UserDataResponse, len(rows))

		for i, row := range rows {
			respObj[i] = UserDataItem{
				ShortURL:    row.ShortURL,
				OriginalURL: row.OriginalURL,
			}
		}

		statusCode := http.StatusOK
		if len(respObj) == 0 {
			statusCode = http.StatusNoContent
		}

		writeResponse(w, respObj, statusCode)
	}
}
