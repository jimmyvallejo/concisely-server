package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sashabaranov/go-openai"
)

func (h *Handlers) ChatGPTCompletion(w http.ResponseWriter, r *http.Request) {
	client := openai.NewClient(h.GPTKEY)

	data, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful assistant that summarizes content clearly and concisely.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: SampleString,
				},
			},
		},
	)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to respond")
		return
	}

	fmt.Printf("Full response: %+v\n", data)

	fmt.Printf("Message content: %s\n", data.Choices[0].Message.Content)
}