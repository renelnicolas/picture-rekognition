package middlewares

import (
	"net/http"
)

// Cors :
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Cache-Control", "max-age=0, no-store, no-cache, must-revalidate")

		if "" != origin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true") // Credentials are cookies, authorization headers or TLS client certificates.
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Type, X-CSRF-Token, Authorization")
			w.Header().Set("Access-Control-Max-Age", "3600")
			return
		}

		next.ServeHTTP(w, r)
	})
}
