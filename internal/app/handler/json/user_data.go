package json

import (
	"github.com/russianlagman/go-musthave-shortener/internal/app/handler"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store"
	"net/http"
)

type UserDataResponse []UserDataItem

type UserDataItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func UserDataHandler(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows := s.ReadUserURLs(handler.ReadContextString(r.Context(), handler.ContextKeyUID{}))
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
