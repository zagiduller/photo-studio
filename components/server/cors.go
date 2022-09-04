package server

import (
	"net/http"
	"strings"
)

// @project photo-studio
// @created 04.09.2022

func AllowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqOrigin = r.Header.Get("Origin")
		if reqOrigin != "" {
			for _, origin := range []string{
				"http://localhost:3000",
				"http://localhost:8081",
			} {
				if origin == reqOrigin {
					header := w.Header()
					header.Add("Access-Control-Allow-Origin", origin)
					header.Add("Access-Control-Allow-Credentials", "true")
					header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS, PUT")
					allowHeaders := []string{
						"Content-Type", "Authorization", "Accept", "Origin",
						"X-Requested-With", "X-Account", "Access-Control-Allow-Headers",
					}
					header.Add("Access-Control-Allow-Headers", strings.Join(allowHeaders, ", "))
					if r.Method == "OPTIONS" {
						w.WriteHeader(http.StatusNoContent)
						return
					}
					break
				}
			}

		}
		h.ServeHTTP(w, r)
	})
}
