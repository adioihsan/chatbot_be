package model

type ChatRequest struct {
	ConversationPID string `json:"conversationPid" `
	Messages        []struct {
		Type    string `json:"type"  validated:"required"`
		Content string `jsong:"content" validated:"required"`
	} `json:"messages"`
}

type ChatResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    *Message `json:"data"`
}
