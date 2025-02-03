package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// func (h *Handlers) AnthropicCompletion(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/event-stream")
// 	w.Header().Set("Transfer-Encoding", "chunked")
// 	w.Header().Set("X-Content-Type-Options", "nosniff")
// 	w.Header().Set("Cache-Control", "no-cache")
// 	w.Header().Set("Connection", "keep-alive")

// 	flusher, ok := w.(http.Flusher)
// 	if !ok {
// 		respondWithError(w, http.StatusInternalServerError, "streaming not supported")
// 		return
// 	}

// 	request := ScrapedDataRequest{}

// 	err := json.NewDecoder(r.Body).Decode(&request)
// 	if err != nil {
// 		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
// 		return
// 	}

// 	formattedContent := formatScrapedContent(request)

// 	client := anthropic.NewClient(request.ApiKey)
// 	resp, err := client.CreateMessagesStream(context.Background(), anthropic.MessagesStreamRequest{
// 		MessagesRequest: anthropic.MessagesRequest{
// 			Model: anthropic.ModelClaude3Haiku20240307,
// 			Messages: []anthropic.Message{
// 				anthropic.NewUserTextMessage("What is your name?"),
// 			},
// 			MaxTokens: 1000,
// 		},
// 		OnContentBlockDelta: func(data anthropic.MessagesEventContentBlockDeltaData) {
// 			fmt.Printf("Stream Content: %s\n", data.Delta.Text)
// 		},
// 	})

// 	req := openai.ChatCompletionRequest{
// 		Model: determineGPTModel(request.Model),
// 		Messages: []openai.ChatCompletionMessage{
// 			{
// 				Role:    openai.ChatMessageRoleSystem,
// 				Content: SystemPrompt,
// 			},
// 			{
// 				Role:    openai.ChatMessageRoleUser,
// 				Content: formattedContent,
// 			},
// 		},
// 		Stream: true,
// 	}

// 	stream, err := client.CreateChatCompletionStream(ctx, req)
// 	if err != nil {
// 		fmt.Printf("ChatCompletionStream error: %v\n", err)
// 		return
// 	}
// 	defer stream.Close()

// 	for {
// 		response, err := stream.Recv()
// 		if errors.Is(err, io.EOF) {
// 			return
// 		}
// 		if err != nil {
// 			fmt.Printf("Stream error: %v\n", err)
// 			fmt.Fprintf(w, "Error: %v", err)
// 			flusher.Flush()
// 			return
// 		}

// 		fmt.Fprint(w, response.Choices[0].Delta.Content)
// 		flusher.Flush()
// 	}
// }

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
