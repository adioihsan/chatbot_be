package handler

import (
	"cms-octo-chat-api/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *BaseHandler) ListMessage(c *fiber.Ctx) error {
	pid, err := uuid.Parse(c.Params("conversation_pid"))
	userId := c.Locals("user_id").(int64)
	limitQ := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitQ)

	if err != nil {
		return fiber.NewError(400, "invalid id")
	}

	var beforePID *uuid.UUID
	if a := c.Query("beforePid"); a != "" {
		if bid, e := uuid.Parse(a); e == nil {
			beforePID = &bid
		}
	}

	msgs, err := h.Repo.ListMessageByConversationPID(c.Context(), userId, pid, limit, beforePID)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.GlobalResponse{
			Code:    fiber.StatusNotFound,
			Message: "Cant find the conversation",
		})
	}
	// get last pid
	var lastPID *uuid.UUID
	if len(msgs) > 0 && len(msgs) >= limit {
		lastPID = &msgs[0].PublicID
	} else {
		lastPID = nil
	}

	return c.Status(fiber.StatusOK).JSON(model.MessageListRes{
		Code:            200,
		Message:         "OK",
		ConversationPID: pid,
		LastPID:         lastPID,
		Data:            &msgs,
	})
}
