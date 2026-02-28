package web

import (
	"context"
	"net/http"
)

type contextKey string

func (app *Application) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := app.Authenticator.ValidateToken(token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var claimsKey contextKey = contextKey(app.Config.App.Slug + "_claims")
		ctx := context.WithValue(r.Context(), claimsKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
