package json

import (
	"github.com/russianlagman/go-musthave-shortener/internal/app/middleware"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store"
	"net/http"
)

type UserDataResponse []UserDataItem

type UserDataItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func UserDataHandler(s store.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows := s.ReadUserURLs(r.Context().Value(middleware.ContextKeyUID{}).(string))
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
