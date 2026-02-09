package jwt

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int32  `json:"id"`
	Salt   string `json:"salt,omitempty"`
	jwt.RegisteredClaims
}

var signingMethod jwt.SigningMethod = jwt.SigningMethodHS256

func Sign(secret string, claims *Claims) (string, error) {
	token := jwt.NewWithClaims(signingMethod, claims)
	return token.SignedString([]byte(secret))
}

func Verify(secret string, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
