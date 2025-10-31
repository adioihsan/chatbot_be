package model

import "github.com/google/uuid"

type SearchResult struct {
	ConversationPid uuid.UUID `json:"conversation_pid"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
}

type SearchResultResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    []SearchResult `json:"data"`
}
