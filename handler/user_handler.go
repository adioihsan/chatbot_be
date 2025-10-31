package handler

import (
	"strconv"

	"cms-octo-chat-api/helper"
	"cms-octo-chat-api/model"

	"github.com/gofiber/fiber/v2"
)

// CreateUser - POST /users
func (h *BaseHandler) CreateUser(c *fiber.Ctx) error {
	body := c.Locals("validatedBody").(*model.UserCreateRequest)

	hashedPassword, err := helper.HashPassword(body.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GlobalResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create user",
			Data:    nil,
		})
	}

	body.Password = hashedPassword

	user, err := helper.Convert[*model.User](body)

	if err := h.Repo.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GlobalResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create user",
			Data:    err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.GlobalResponse{
		Code:    fiber.StatusCreated,
		Message: "User created successfully",
		Data:    body,
	})
}

// GetUserByID - GET /users/:id
func (h *BaseHandler) GetUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GlobalResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid user ID",
			Data:    nil,
		})
	}

	user, err := h.Repo.GetUserByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.GlobalResponse{
			Code:    fiber.StatusNotFound,
			Message: "User not found",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GlobalResponse{
		Code:    fiber.StatusOK,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// UpdateUser - PUT /users/:id
func (h *BaseHandler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GlobalResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid user ID",
			Data:    nil,
		})
	}

	body := c.Locals("validatedBody").(*model.User)
	body.ID = int64(id) // ensure the ID is set from URL param

	if err := h.Repo.UpdateUser(body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GlobalResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update user",
			Data:    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GlobalResponse{
		Code:    fiber.StatusOK,
		Message: "User updated successfully",
		Data:    body,
	})
}

// DeleteUser - DELETE /users/:id
func (h *BaseHandler) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GlobalResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid user ID",
			Data:    nil,
		})
	}

	if err := h.Repo.DeleteUser(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GlobalResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete user",
			Data:    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.GlobalResponse{
		Code:    fiber.StatusOK,
		Message: "User deleted successfully",
		Data:    nil,
	})
}

// Create user matrix
func (h *BaseHandler) CreateUserMatrix(c *fiber.Ctx) error {

	body := c.Locals("validatedBody").(*model.UserMatrixRequest)

	userId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.GlobalResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid user ID",
			Data:    nil,
		})
	}

	user, err := h.Repo.GetUserByID(uint(userId))

	userMatrix, err := helper.Convert(body, func(u *model.UserMatrix) error {
		u.UserID = user.ID
		return nil
	})

	if err := h.Repo.CreateUserMatrix(&userMatrix); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GlobalResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create user matrix",
			Data:    err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.GlobalResponse{
		Code:    fiber.StatusCreated,
		Message: "User matrix created or updated successfully",
		Data:    body,
	})
}

// user short data
func (h *BaseHandler) Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	userPID := c.Locals("user_pid").(string)
	userEmail := c.Locals("email").(string)
	userName := c.Locals("user_name").(string)

	var userMe model.UserMe
	userMe.Name = userName
	userMe.PublicID = userPID
	userMe.Email = userEmail

	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(model.GlobalResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "You need to login"})
	}
	return c.Status(fiber.StatusOK).JSON(model.UserMeResponse{
		Code:    fiber.StatusOK,
		Message: "Authorized",
		Data:    userMe,
	})
}
