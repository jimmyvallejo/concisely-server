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
	case gpt4Point1:
		return openai.ChatModelGPT4_1
	case gptReasoning:
		return openai.ChatModelO4Mini
	default:
		return openai.ChatModelGPT4oMini
	}
}

func determineAnthropicModel(model string) anthropic.Model {
	switch model {
	case claude3Point5Haiku:
		return anthropic.ModelClaude3_5HaikuLatest
	case claude4Opus:
		return anthropic.ModelClaudeOpus4_0
	case claude4Sonnet:
		return anthropic.ModelClaudeSonnet4_0
	default:
		return anthropic.ModelClaudeSonnet4_0
	}
}
