package handler

import (
	"cms-octo-chat-api/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *BaseHandler) GlobalSearch(c *fiber.Ctx) error {
	ctx := c.UserContext()
	userID := c.Locals("user_id").(int64)
	limitQ := c.Query("limit", "10")
	limit, _ := strconv.Atoi(limitQ)
	searchQuery := c.Query("query")

	if len(searchQuery) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(model.GlobalResponse{
			Code:    403,
			Message: "Search query cant be empty !",
		})
	}
	// searchQuery = helper.Prefixify(searchQuery)

	searchResult, err := h.Repo.FindChats(ctx, searchQuery, limit, userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.GlobalResponse{
			Code:    500,
			Message: "Server Error !",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.SearchResultResponse{
		Code:    200,
		Message: "Search result",
		Data:    searchResult,
	})

}
