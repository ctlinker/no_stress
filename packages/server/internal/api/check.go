package api

import (
	"encoding/json"
	"log"
	"net/http"
	"server/internal/env"
	"server/internal/jwt"
)

func Check(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := (ctx.Value("user")).(*jwt.Claims)
	if !ok {
		http.Error(w, "You're an Anon User", http.StatusUnauthorized)
		return
	}

	if env.Load().RUNNING_IN == string(env.TEST) {
		jsonf, err := json.Marshal(jwt.DebugClaims{
			Sub: claims.Subject,
			Exp: claims.ExpiresAt.Unix(),
		})
		if err != nil {
			log.Printf("failed to marshal JWT claims for logging: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("jwt claims: %s", jsonf)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("You're Authenicated"))
}
