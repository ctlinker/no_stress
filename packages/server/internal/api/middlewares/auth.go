package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"server/internal/env"
	"server/internal/jwt"
)

func WithAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := jwt.VerifyJWT(env.Load().JWT_ACCESS_SECRET, c.Value)
		if err != nil {
			if env.Load().RUNNING_IN == string(env.TEST) {
				log.Printf("%v", fmt.Errorf("[middlewares] error while check access token : %v", err))
			}
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
