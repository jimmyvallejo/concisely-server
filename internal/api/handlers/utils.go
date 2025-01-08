package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
