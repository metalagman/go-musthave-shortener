package json

import (
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
		respObj := UserDataResponse{}
		writeResponse(w, respObj, http.StatusOK)
	}
}
