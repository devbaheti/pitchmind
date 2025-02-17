package llm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Message Message `json:"message"`
}

type LLM struct {
	Client http.Client
}

func New() LLM {
	return LLM{
		Client: http.Client{Timeout: 30 * time.Second},
	}
}

func (l LLM) GenerateClueFromLLM(clues []string, currentPlayer string, prompt string) string {
	messages := []Message{
		{
			Role: "system",
			Content: fmt.Sprintf(
				"You're a football expert comedian. The player is %s. "+
					"Give %d funny, witty clues using wordplay, pop culture refs, and funny analogies. "+
					"Never mention the name, team, or nationality directly. "+
					"Make it humorous and engaging under 40 characters! Format: emoji + clue",
				currentPlayer, len(clues)+1,
			),
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	requestBody := ChatRequest{
		Model:    "llama3.2",
		Messages: messages,
		Stream:   false,
	}

	jsonBody, _ := json.Marshal(requestBody)
	resp, err := l.Client.Post("http://localhost:11434/api/chat", "application/json", strings.NewReader(string(jsonBody)))
	if err != nil {
		return "🚧 Oops, the clue machine broke! Try another guess..."
	}
	defer resp.Body.Close()

	var chatResp ChatResponse
	json.NewDecoder(resp.Body).Decode(&chatResp)
	return chatResp.Message.Content
}
