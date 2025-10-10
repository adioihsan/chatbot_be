package middleware

import (
	"cms-octo-chat-api/helper"
	"cms-octo-chat-api/model"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Claims defines the structure of JWT claims

// AuthMiddleware checks for valid JWT in Authorization header
func (m *BaseMiddleware) JwtAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.GlobalResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Missing Authorization header",
				Data:    nil,
			})
		}

		// Expect format: Bearer token_here
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.GlobalResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Invalid Authorization header format",
				Data:    nil,
			})
		}

		tokenString := tokenParts[1]

		// Parse token

		claims, err := helper.ValidateJWT(tokenString, m.Env.JWTSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(model.GlobalResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Invalid or expired token",
				Data:    nil,
			})
		}

		// Check expiration
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.GlobalResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Token has expired",
				Data:    nil,
			})
		}

		// Store user data in context for later use
		c.Locals("user_id", claims.UserID)
		c.Locals("user_pid", claims.UserPID)
		c.Locals("email", claims.Email)
		c.Locals("user_matrix", claims.UserMatrix)

		return c.Next()
	}
}
