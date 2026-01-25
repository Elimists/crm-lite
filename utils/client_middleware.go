package utils

import (
	"context"
	"log"
	"net/http"
)

func ClientMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		source := r.Header.Get("Origin")
		if source == "" {
			log.Println("incoming request is missing origin header")
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		client, ok := clientsByOrigin[source]
		if !ok {
			log.Println("origin not allowed:", source)
			http.Error(w, "not allowed", http.StatusForbidden)
			return
		}

		// call the next handler
		ctx := context.WithValue(r.Context(), clientContextKey, client)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
