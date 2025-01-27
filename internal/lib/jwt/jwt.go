package jwt

import (
	"fmt"
	"time"

	"github.com/Melikhov-p/sso-grpc-go/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int64
	AppID  int32
}

func BuildJWTToken(userID int64, app *models.App, tokenLifeTime time.Duration) (string, error) {
	op := "lib.jwt.BuildJWTToken"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID: userID,
		AppID:  app.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifeTime)),
		},
	})

	tokenString, err := token.SignedString([]byte(app.SecretKey))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, nil
}

func GetUserIDFromToken(tokenString string, secretKey string) (int64, error) {
	op := "lib.jwt.GetUserIDFromToken"
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method for token")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return -1, fmt.Errorf("%s: error parsing token: %w", op, err)
	}
	if !token.Valid {
		return -1, fmt.Errorf("%s: invalid token", op)
	}

	return claims.UserID, nil
}
