package jwt

import (
	"context"
	"errors"
	"server/internal/db"
	"server/schema"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenAccessToken(access_secret string, userID int32) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "ACCESS",
		},
	}

	return Sign(access_secret, claims)
}

func GenRefreshToken(refresh_secret string, userID int32, db *db.DB, ctx context.Context) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "REFRESH",
		},
	}

	token, err := Sign(refresh_secret, claims)
	if err != nil {
		return "", err
	}

	_, err = db.CreateSession(ctx, schema.CreateSessionParams{
		UserID:    userID,
		UpdatedAt: time.Now(),
		TokenHash: HashToken(token),
	})

	if err != nil {
		return "", errors.New("session create failed")
	}

	return token, nil
}
