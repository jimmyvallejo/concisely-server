package handlers

import (
	"context"
	"encoding/json"
	// "encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/sashabaranov/go-openai"
)

func (h *Handlers) ChatGPTCompletion(w http.ResponseWriter, r *http.Request) {
	client := openai.NewClient(h.GPTKEY)
	ctx := context.Background()

	request := ScrapedDataRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	formattedContent := formatScrapedContent(request)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: SystemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: formattedContent,
			},
		},
		Stream: true,
	}

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Printf("Stream response: ")

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		fmt.Printf("%s", response.Choices[0].Delta.Content)
	}
}
