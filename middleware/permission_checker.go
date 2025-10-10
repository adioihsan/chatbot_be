package middleware

import (
	"strings"

	"cms-octo-chat-api/model"

	"github.com/gofiber/fiber/v2"
)

func (m *BaseMiddleware) PermissionChecker(permission string) fiber.Handler {
	required := strings.ToUpper(strings.TrimSpace(permission))
	if required == "" {
		required = "R" // default to Read
	}

	// map string -> enum
	var action model.PermissionActionEnum
	switch required {
	case "C":
		action = model.Create
	case "R":
		action = model.Read
	case "U":
		action = model.Update
	case "D":
		action = model.Delete
	case "A":
		action = model.Upload
	case "B":
		action = model.Download
	default:
		// unknown permission code
		return func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusBadRequest).JSON(model.GlobalResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Unknown permission code",
				Data:    nil,
			})
		}
	}

	return func(c *fiber.Ctx) error {
		// 1) Require user_id
		if c.Locals("user_id") == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(model.GlobalResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Not Authorized",
				Data:    nil,
			})
		}

		// 2) Get matrix (prefer set by previous auth middleware)
		lm := c.Locals("user_matrix")
		if lm == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(model.GlobalResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Not Authorized",
				Data:    nil,
			})
		}

		// 3) Type assert safely
		var matrix *model.UserMatrix
		switch v := lm.(type) {
		case *model.UserMatrix:
			matrix = v
		case model.UserMatrix:
			matrix = &v
		default:
			return c.Status(fiber.StatusUnauthorized).JSON(model.GlobalResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Invalid permission matrix",
				Data:    nil,
			})
		}

		// 4) Check permission
		if !hasPermission(matrix, action) {
			return c.Status(fiber.StatusForbidden).JSON(model.GlobalResponse{
				Code:    fiber.StatusForbidden,
				Message: "Permission Denied",
				Data:    nil,
			})
		}

		// 5) Pass along (optional: re-store normalized matrix)
		c.Locals("user_matrix", matrix)
		return c.Next()
	}
}

// Correct signature: returns bool
func hasPermission(matrix *model.UserMatrix, permission model.PermissionActionEnum) bool {
	switch permission {
	case model.Create:
		return ptrTrue(matrix.IsCreate)
	case model.Read:
		return ptrTrue(matrix.IsRead)
	case model.Update:
		return ptrTrue(matrix.IsUpdate)
	case model.Delete:
		return ptrTrue(matrix.IsDelete)
	case model.Upload:
		return ptrTrue(matrix.IsUpload)
	case model.Download:
		return ptrTrue(matrix.IsDownload)
	case model.Archive:
		return ptrTrue(matrix.IsArchive)
	default:
		return false
	}
}

// helper: true only if non-nil and true
func ptrTrue(b *bool) bool { return b != nil && *b }
