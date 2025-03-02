package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/openai/openai-go"
)

func UnmarshalCompletionRequest(data []byte) (CompletionRequest, error) {
	var requestType struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(data, &requestType); err != nil {
		return nil, fmt.Errorf("failed to parse request type: %w", err)
	}

	switch requestType.Type {
	case "web":
		var req ScrapedDataRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, fmt.Errorf("failed to parse scraped data request: %w", err)
		}
		return &req, nil

	case "pdf":
		var req ScrapedDataRequestPDF
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, fmt.Errorf("failed to parse prompt request: %w", err)
		}
		return &req, nil

	default:
		return nil, fmt.Errorf("unknown request type: %s", requestType.Type)
	}
}

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

func formatScrapedContent(req ScrapedDataRequest) string {
	var formattedContent strings.Builder

	formattedContent.WriteString(fmt.Sprintf("Title: %s\n\n", req.Title))

	if len(req.Headers) > 0 {
		formattedContent.WriteString("Headers:\n")
		for _, header := range req.Headers {
			formattedContent.WriteString(fmt.Sprintf("- %s: %s\n", header.Type, header.Text))
		}
		formattedContent.WriteString("\n")
	}

	if req.MetaDescription != nil && *req.MetaDescription != "" {
		formattedContent.WriteString(fmt.Sprintf("Meta Description: %s\n\n", *req.MetaDescription))
	}

	if req.MainContent != nil && *req.MainContent != "" {
		formattedContent.WriteString(fmt.Sprintf("Main Content:\n%s\n\n", *req.MainContent))
	}

	if len(req.Paragraphs) > 0 {
		formattedContent.WriteString("Paragraphs:\n")
		for _, para := range req.Paragraphs {
			formattedContent.WriteString(fmt.Sprintf("%s\n\n", para))
		}
	}

	if len(req.Links) > 0 {
		formattedContent.WriteString("Relevant Links:\n")
		for _, link := range req.Links {
			formattedContent.WriteString(fmt.Sprintf("- %s: %s\n", link.Text, link.Href))
		}
	}

	return formattedContent.String()
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
	case Claude3Haiku:
		return anthropic.ModelClaude3_5HaikuLatest
	case Claude3Opus:
		return anthropic.ModelClaude3OpusLatest
	case Claude3Point7Sonnet:
		return anthropic.ModelClaude3_7SonnetLatest
	default:
		return anthropic.ModelClaude3_7SonnetLatest
	}
}
