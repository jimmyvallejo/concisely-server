package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/openai/openai-go"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Error: message})
}

func respondNoBody(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
}

func determineGPTModel(model string) string {
	switch model {
	case gpt4oMini:
		return openai.ChatModelGPT4oMini
	case gpt4oStandard:
		return openai.ChatModelGPT4o
	case gpt4Old:
		return openai.ChatModelGPT4Turbo
	default:
		return openai.ChatModelGPT4oMini
	}
}

func determineAnthropicModel(model string) string {
	switch model {
	case claude3Haiku:
		return anthropic.ModelClaude3_5HaikuLatest
	case claude3Opus:
		return anthropic.ModelClaude3OpusLatest
	case claude3Point7Sonnet:
		return anthropic.ModelClaude3_7SonnetLatest
	default:
		return anthropic.ModelClaude3_7SonnetLatest
	}
}
