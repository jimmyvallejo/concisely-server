package handlers

import (
	"encoding/json"
	"time"

	"fmt"
	"net/http"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
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

	client := openai.NewClient(option.WithAPIKey(request.ApiKey))

	stream := client.Chat.Completions.NewStreaming(
		r.Context(),
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(SystemPrompt),
				openai.UserMessage(formattedContent),
			}),
			Model: openai.F(determineGPTModel(request.Model)),
		},
	)

	acc := openai.ChatCompletionAccumulator{}

	rateLimiter := time.NewTicker(25 * time.Millisecond)
	defer rateLimiter.Stop()

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			<-rateLimiter.C
			fmt.Fprint(w, chunk.Choices[0].Delta.Content)
			flusher.Flush()
		}
	}

	if err := stream.Err(); err != nil {
		fmt.Printf("Stream error: %v\n", err)
		fmt.Fprintf(w, "Error: %v", err)
		flusher.Flush()
		return
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
