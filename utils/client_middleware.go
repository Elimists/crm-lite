package utils

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
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

		signature := r.Header.Get("X-Signature")
		if signature == "" {
			log.Println("missing HMAC signature")
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		log.Printf("hmac signature recieved:%s\n", signature)
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("failed to read body: ", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		log.Printf("body string: %s\n", string(bodyBytes))

		secret := os.Getenv("SHARED_API_HMAC_KEY")
		if !VerifyHMAC(bodyBytes, signature, secret) {
			log.Println("invalid HMAC signature")
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		// call the next handler
		ctx := context.WithValue(r.Context(), clientContextKey, client)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
