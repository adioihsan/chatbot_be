package handler

import (
	"cms-octo-chat-api/helper"
	"cms-octo-chat-api/model"

	"github.com/gofiber/fiber/v2"
)

// CreateUser - POST /users
func (h *BaseHandler) Login(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(*model.AuthRequest)

	user, err := h.Repo.GetUserByEmail(body.Email)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.GlobalResponse{
			Code:    fiber.StatusNotFound,
			Message: "User not found",
			Data:    nil,
		})
	}

	if err := helper.CheckPassword(body.Password, user.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorValidationResponse{
			Code:    fiber.ErrBadRequest.Code,
			Message: "Wrong Password",
			Errors: map[string]string{
				"password": "The password must be at least 8 characters",
			},
		})
	}

	token, err := helper.GenerateJWT(user, h.Env.JWTSecret)

	if err != nil {
		h.Logs.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.AuthSuccessResponse{
		Code:    fiber.StatusOK,
		Message: "Login Success",
		Token:   token,
	})

}
