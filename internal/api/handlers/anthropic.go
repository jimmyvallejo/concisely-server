package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func (h *Handlers) AnthropicCompletion(w http.ResponseWriter, r *http.Request) {
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	var request ScrapedDataRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	apiKey := request.GetAPIKey()
	model := request.GetModel()
	formattedContent := request.FormatContent()

	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)

	stream := client.Messages.NewStreaming(
		r.Context(),
		anthropic.MessageNewParams{
			Model:     anthropic.F(determineAnthropicModel(model)),
			MaxTokens: anthropic.F(int64(2024)),
			System: anthropic.F([]anthropic.TextBlockParam{
				anthropic.NewTextBlock(SystemPrompt),
			}),
			Messages: anthropic.F([]anthropic.MessageParam{
				anthropic.NewUserMessage(anthropic.NewTextBlock(formattedContent)),
			}),
		},
	)

	rateLimiter := time.NewTicker(10 * time.Millisecond)
	defer rateLimiter.Stop()

	message := anthropic.Message{}

	for stream.Next() {
		event := stream.Current()
		message.Accumulate(event)

		switch delta := event.Delta.(type) {
		case anthropic.ContentBlockDeltaEventDelta:
			if delta.Text != "" {
				<-rateLimiter.C
				fmt.Fprint(w, delta.Text)
				flusher.Flush()
			}
		}
	}

	if err := stream.Err(); err != nil {
		var apiErr *anthropic.Error
		if errors.As(err, &apiErr) {
			respondWithError(w, apiErr.StatusCode, fmt.Sprintf("Stream error: %v", apiErr))
		} else {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Stream error: %v", err))
		}
		return
	}
}

func (h *Handlers) ValidateAnthropicKey(w http.ResponseWriter, r *http.Request) {
	request := ValidateKeyRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	client := anthropic.NewClient(
		option.WithAPIKey(request.APIKey),
	)

	_, err = client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude3_5HaikuLatest),
		MaxTokens: anthropic.F(int64(1)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What is a quaternion?")),
		}),
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	respondNoBody(w, http.StatusOK)
}
