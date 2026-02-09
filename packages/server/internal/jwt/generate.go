package jwt

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
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

func genSalt(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func GenRefreshToken(refresh_secret string, userID int32, db *db.DB, ctx context.Context) (string, error) {

	salt, err := genSalt(16) // 16 bytes = 32 hex chars
	if err != nil {
		return "", fmt.Errorf("cannot generate salt: %w", err)
	}

	claims := &Claims{
		UserID: userID,
		Salt:   salt,
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
		return "", err
	}

	return token, nil
}
