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
		for _, row := range rows {
			respObj = append(respObj, UserDataItem{
				ShortURL:    row.ShortURL,
				OriginalURL: row.OriginalURL,
			})
		}

		writeResponse(w, respObj, http.StatusOK)
	}
}
