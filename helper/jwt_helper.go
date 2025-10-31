package helper

import (
	"cms-octo-chat-api/model"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT generates a JWT token for a given user
func GenerateJWT(user *model.User, jwtSecret string) (string, error) {
	claims := model.JWTClaims{
		UserID:     user.ID,
		UserPID:    user.PublicID.String(),
		Email:      user.Email,
		Name:       user.Name,
		UserMatrix: user.UserMatrix,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)), // 3 day expiry
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "cms-chat",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}

// ValidateJWT parses and validates a JWT token
func ValidateJWT(tokenString, jwtSecret string) (*model.JWTClaims, error) {
	claims := &model.JWTClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(jwtSecret), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
