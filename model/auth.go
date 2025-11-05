package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email,exists=users=email"`
	Password string `json:"password" validate:"required,min=8"`
}

type AuthSuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type AuthWithTokenRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Token    string `json:"token" validated:"required"`
}

type JWTClaims struct {
	UserID               int64  `json:"UserID"`  // matches "UserID"
	UserPID              string `json:"UserPID"` // matches "UserPID" (UUID as string)
	Email                string `json:"Email"`   // matches "Email"
	Name                 string `json:"name"`
	UserMatrix           any    `json:"UserMatrix"` // or a concrete type if you have one
	jwt.RegisteredClaims        // covers iss, exp, iat, etc.
}
