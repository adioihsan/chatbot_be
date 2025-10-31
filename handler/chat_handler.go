package handler

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"cms-octo-chat-api/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openai "github.com/openai/openai-go/v3"
)

func (h *BaseHandler) SingleChat(c *fiber.Ctx) error {
	ctx := c.UserContext()

	if ctx == nil {
		ctx = context.Background()
	}
	body := c.Locals("validatedBody").(*model.ChatRequest)
	userID := c.Locals("user_id").(int64)
	cid := body.ConversationPID

	convoPID, err := uuid.Parse(cid)
	if err != nil {
		return fiber.NewError(400, "invalid conversationId")
	}

	convo, err := h.Repo.GetConversationByPublicID(ctx, convoPID, userID)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	var prevMessageId int64

	for _, m := range body.Messages {
		if m.Type != "text" {
			continue
		}

		var refId *int64
		if prevMessageId > 0 {
			ref := prevMessageId
			refId = &ref
		}

		nm, _ := h.Repo.AppendMessage(ctx, &model.Message{
			ConversationID: convo.ID,
			Role:           "user",
			Content:        m.Content,
			RefID:          refId,
		})
		_ = h.Repo.TouchConversation(ctx, convo.ID)
		_ = h.Repo.RebuildMessageContentFTS(ctx, nm.ID)

		if prevMessageId == 0 {
			prevMessageId = nm.ID
		}
	}

	messagesHistory, err := h.Repo.ListMessageByConversationPID(ctx, userID, convo.PublicID, 50, nil)

	params := openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(h.Env.OpenAiModel),
		Messages: []openai.ChatCompletionMessageParamUnion{},
	}

	params.Messages = append(params.Messages, openai.SystemMessage("You are ChatGPT, an AI assistant that always responds in GitHub-Flavored Markdown (GFM) format. Use code blocks for code, bullet points for lists, and headings where appropriate. Always format responses cleanly for display in Markdown renderers."))

	for _, mh := range messagesHistory {
		switch mh.Role {
		case "assistant":
			params.Messages = append(params.Messages, openai.UserMessage(mh.Content))
		// case "system:
		// 	params.Messages = append(params.Messages,openai.SystemMessage(mh.Content))
		default:
			params.Messages = append(params.Messages, openai.UserMessage(mh.Content))
		}
	}

	for _, m := range body.Messages {
		params.Messages = append(params.Messages, openai.UserMessage(m.Content))
	}

	chatResponse, err := h.OpenAI.Chat.Completions.New(ctx, params)
	chatContent := chatResponse.Choices[0].Message.Content

	if err != nil {
		return fiber.NewError(400, "System Error")
	}

	message, err := h.Repo.AppendMessage(ctx, &model.Message{
		ConversationID: convo.ID,
		Role:           model.ROLE_ASSISTANT,
		Content:        chatContent,
		RefID:          &prevMessageId,
	})

	err = h.Repo.RebuildMessageContentFTS(ctx, message.ID)

	return c.Status(fiber.StatusOK).JSON(model.ChatResponse{
		Code:    fiber.StatusOK,
		Message: "OK",
		Data:    message,
	})
}

func (h *BaseHandler) StreamChat(c *fiber.Ctx) error {

	// / headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("X-Accel-Buffering", "no") // avoid proxy buffering

	ctx := c.UserContext()

	if ctx == nil {
		ctx = context.Background()
	}

	body := c.Locals("validatedBody").(*model.ChatRequest)

	userID := c.Locals("user_id").(int64)

	cid := body.ConversationPID

	convoPID, err := uuid.Parse(cid)
	if err != nil {
		return fiber.NewError(400, "invalid conversationId")
	}

	convo, err := h.Repo.GetConversationByPublicID(ctx, convoPID, userID)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	var prevMessageId int64

	for _, m := range body.Messages {
		if m.Type != "text" {
			continue
		}

		var refId *int64
		if prevMessageId > 0 {
			ref := prevMessageId
			refId = &ref
		}

		nm, _ := h.Repo.AppendMessage(ctx, &model.Message{
			ConversationID: convo.ID,
			Role:           "user",
			Content:        m.Content,
			RefID:          refId,
		})
		_ = h.Repo.TouchConversation(ctx, convo.ID)

		if prevMessageId == 0 {
			prevMessageId = nm.ID
		}
	}

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer w.Flush()

		fmt.Fprint(w, ":ok\n\n")
		w.Flush()

		params := openai.ChatCompletionNewParams{
			Model:    openai.ChatModel(h.Env.OpenAiModel),
			Messages: []openai.ChatCompletionMessageParamUnion{},
		}
		for _, m := range body.Messages {
			params.Messages = append(params.Messages, openai.UserMessage(m.Content))
		}

		stream := h.OpenAI.Chat.Completions.NewStreaming(ctx, params)
		defer stream.Close()

		var sb strings.Builder
		var gotAny bool

		for stream.Next() {
			chunk := stream.Current()
			if len(chunk.Choices) > 0 {
				if delta := chunk.Choices[0].Delta.Content; delta != "" {
					gotAny = true
					sb.WriteString(delta)
					fmt.Fprintf(w, "data: {\"type\":\"text-delta\",\"textDelta\":%q}\n\n", delta)
					w.Flush()
				}
			}
		}

		if err := stream.Err(); err != nil {
			fmt.Fprintf(w, "data: {\"type\":\"error\",\"message\":%q}\n\n", err.Error())
			w.Flush()
			return
		}

		assistantText := sb.String()

		if gotAny && strings.TrimSpace(assistantText) != "" {
			_, _ = h.Repo.AppendMessage(ctx, &model.Message{
				ConversationID: convo.ID,
				Role:           "assistant",
				Content:        assistantText,
				RefID:          &prevMessageId,
			})
			_ = h.Repo.TouchConversation(ctx, convo.ID)
		}

		fmt.Fprintf(w, "data: {\"type\":\"done\",\"conversationId\":%q}\n\n", convo.PublicID.String())
		w.Flush()
	})
	return nil

}
