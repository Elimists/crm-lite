package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

var allowedOrigins map[string]bool

func LoadAllowedOrigins(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var origins []string
	err = json.Unmarshal(data, &origins)
	if err != nil {
		return err
	}

	allowedOrigins = make(map[string]bool)
	for _, o := range origins {
		allowedOrigins[o] = true
	}

	return nil
}

func OriginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		source := r.Header.Get("Origin")
		if source == "" {
			log.Println("incoming request is missing origin header")
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if !allowedOrigins[source] {
			log.Println("origin not allowed:", source)
			http.Error(w, "not allowed", http.StatusForbidden)
			return
		}

		// call the next handler
		next.ServeHTTP(w, r)
	})
}
