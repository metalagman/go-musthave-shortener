package middleware

import (
	"net/http"
)

func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//if r.Header.Get("Content-Type") != "application/json" {
		//	http.Error(w, "invalid content type", http.StatusBadRequest)
		//	return
		//}
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
