package handler

import (
	"cms-octo-chat-api/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *BaseHandler) ListConversation(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	limitQ := c.Query("limit", "50")
	limit, err := strconv.Atoi(limitQ)

	var beforePID *uuid.UUID
	if a := c.Query("beforePid"); a != "" {
		if bid, e := uuid.Parse(a); e == nil {
			beforePID = &bid
		}
	}

	conversations, err := h.Repo.ListConversation(c.Context(), userID, limit, beforePID)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	// get last pid
	var lastPID *uuid.UUID
	if len(conversations) > 0 && len(conversations) >= limit {
		lastPID = &conversations[len(conversations)-1].PublicID
	} else {
		lastPID = nil
	}

	return c.Status(fiber.StatusOK).JSON(model.ConversationListRes{
		Code:    200,
		LastPid: lastPID,
		Message: "Conversation List",
		Data:    conversations,
	})
}

func (h *BaseHandler) CreateConversation(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	body := c.Locals("validatedBody").(*model.ConversationCreateReq)
	conversation, err := h.Repo.CreateConversation(c.Context(), body.Title, userID)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}
	return c.Status(201).JSON(model.ConversationCreateRes{
		Code:    201,
		Message: "Conversation created successfully",
		Data:    conversation,
	})
}

// func (h *BaseHandler) getConversation(c *fiber.Ctx) error {
// 	pid, err := uuid.Parse(c.Params("pid"))
// 	if err != nil {
// 		return fiber.NewError(400, "invalid id")
// 	}
// 	cv, err := h.Repo.GetConversationByPublicID(c.Context(), pid, c.Get("user_id"))
// 	if err != nil {
// 		return fiber.ErrNotFound
// 	}
// 	return c.JSON(cv)
// }

func (h *BaseHandler) RenameConversation(c *fiber.Ctx) error {
	pid, err := uuid.Parse(c.Params("conversation_pid"))
	userID := c.Locals("user_id").(int64)

	if err != nil {
		return fiber.NewError(400, "invalid id")
	}

	body := c.Locals("validatedBody").(*model.ConversationRenameReq)

	if err := h.Repo.RenameConversation(c.Context(), pid, userID, body.Title); err != nil {
		return fiber.NewError(500, err.Error())
	}
	return c.Status(202).JSON(model.GlobalResponse{
		Code:    202,
		Message: "Conversation renamed successfully",
	})
}

func (h *BaseHandler) RemoveConversation(c *fiber.Ctx) error {
	pid, err := uuid.Parse(c.Params("pid"))
	userID := c.Locals("user_id").(int64)
	if err != nil {
		return fiber.NewError(400, "invalid id")
	}
	if err := h.Repo.DeleteConversation(c.Context(), pid, userID); err != nil {
		return fiber.NewError(500, err.Error())
	}
	return c.SendStatus(204)
}
