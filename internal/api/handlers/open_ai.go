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
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "streaming not supported")
		return
	}

	request := ScrapedDataRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	formattedContent := formatScrapedContent(request)

	client := openai.NewClient(h.GPTKEY)
	ctx := context.Background()

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

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			fmt.Fprintf(w, "Error: %v", err)
			flusher.Flush()
			return
		}

		fmt.Fprint(w, response.Choices[0].Delta.Content)
		flusher.Flush()
	}
}

func (h *Handlers) ValidateOpenAIKey(w http.ResponseWriter, r *http.Request) {
	request := ValidateKeyRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	req, err := http.NewRequest("GET", "https://api.openai.com/v1/models", nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create request")
		return
	}

	req.Header.Add("Authorization", "Bearer "+request.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to make request")
		return
	}
	defer resp.Body.Close()


	if resp.StatusCode != http.StatusOK {
		respondWithError(w, resp.StatusCode, "Invalid API key")
		return
	}

	respondNoBody(w, http.StatusOK)
}
