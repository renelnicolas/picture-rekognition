package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Logger :
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if "/favicon.ico" != r.URL.String() {
			log.Println(strings.Repeat(">", 15) + " Logger " + strings.Repeat("<", 15))
			log.Printf("%s || %s\n", r.Host, r.URL.Path)

			for k, v := range r.Header {
				fmt.Printf("%v : %v\n", strings.ToLower(k), string(v[0]))
			}
		}

		next.ServeHTTP(w, r)
	})
}
