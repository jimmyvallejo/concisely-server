package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	var request ScrapedDataRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to read request body: "+err.Error())
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
			Model:     determineAnthropicModel(model),
			MaxTokens: int64(2024),
			System: []anthropic.TextBlockParam{
				{Text: systemPromptWeb},
			},
			Messages: []anthropic.MessageParam{
				anthropic.NewUserMessage(anthropic.NewTextBlock(formattedContent)),
			},
		},
	)

	rateLimiter := time.NewTicker(10 * time.Millisecond)
	defer rateLimiter.Stop()

	message := anthropic.Message{}

	for stream.Next() {
		event := stream.Current()
		message.Accumulate(event)

		switch eventVariant := event.AsAny().(type) {
		case anthropic.ContentBlockDeltaEvent:
			switch deltaVariant := eventVariant.Delta.AsAny().(type) {
			case anthropic.TextDelta:
				if deltaVariant.Text != "" {
					<-rateLimiter.C
					fmt.Fprint(w, deltaVariant.Text)
					flusher.Flush()
				}
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
		respondWithError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	client := anthropic.NewClient(
		option.WithAPIKey(request.APIKey),
	)
	_, err = client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfText: &anthropic.TextBlockParam{Text: "What is a quaternion?"},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaude3_7SonnetLatest,
	})
	if err != nil {
		panic(err.Error())
	}

	respondNoBody(w, http.StatusOK)
}
