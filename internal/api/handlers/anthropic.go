package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/liushuangls/go-anthropic/v2"
)

func (h *Handlers) ValidateAnthropicKey(w http.ResponseWriter, r *http.Request) {
	request := ValidateKeyRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	client := anthropic.NewClient(request.APIKey)

	_, err = client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model: anthropic.ModelClaude3Haiku20240307,
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage("Hello"),
		},
		MaxTokens: 1,
	})

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	respondNoBody(w, http.StatusOK)
}
