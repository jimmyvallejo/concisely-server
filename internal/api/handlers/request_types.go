package handlers

import (
	"fmt"
	"strings"
)

type CompletionRequest interface {
	GetAPIKey() string
	GetModel() string
	FormatContent() string
}

type ScrapedDataRequest struct {
	Title           string   `json:"title"`
	Headers         []Header `json:"headers"`
	Paragraphs      []string `json:"paragraphs"`
	Links           []Link   `json:"links"`
	MetaDescription *string  `json:"metaDescription"`
	MainContent     *string  `json:"mainContent"`
	ApiKey          string   `json:"apiKey"`
	Model           string   `json:"model"`
	Type            string   `json:"type"`
}

func (r *ScrapedDataRequest) GetAPIKey() string {
	return r.ApiKey
}

func (r *ScrapedDataRequest) GetModel() string {
	return r.Model
}

type ScrapedDataRequestPDF struct {
	Title   string  `json:"title"`
	Content *string `json:"mainContent"`
	ApiKey  string  `json:"apiKey"`
	Model   string  `json:"model"`
	Type    string  `json:"type"`
}

func (r *ScrapedDataRequestPDF) GetAPIKey() string {
	return r.ApiKey
}

func (r *ScrapedDataRequestPDF) GetModel() string {
	return r.Model
}

func (r *ScrapedDataRequest) FormatContent() string {
	var formattedContent strings.Builder

	formattedContent.WriteString(fmt.Sprintf("Title: %s\n\n", r.Title))

	if len(r.Headers) > 0 {
		formattedContent.WriteString("Headers:\n")
		for _, header := range r.Headers {
			formattedContent.WriteString(fmt.Sprintf("- %s: %s\n", header.Type, header.Text))
		}
		formattedContent.WriteString("\n")
	}

	if r.MetaDescription != nil && *r.MetaDescription != "" {
		formattedContent.WriteString(fmt.Sprintf("Meta Description: %s\n\n", *r.MetaDescription))
	}

	if r.MainContent != nil && *r.MainContent != "" {
		formattedContent.WriteString(fmt.Sprintf("Main Content:\n%s\n\n", *r.MainContent))
	}

	if len(r.Paragraphs) > 0 {
		formattedContent.WriteString("Paragraphs:\n")
		for _, para := range r.Paragraphs {
			formattedContent.WriteString(fmt.Sprintf("%s\n\n", para))
		}
	}

	if len(r.Links) > 0 {
		formattedContent.WriteString("Relevant Links:\n")
		for _, link := range r.Links {
			formattedContent.WriteString(fmt.Sprintf("- %s: %s\n", link.Text, link.Href))
		}
	}

	return formattedContent.String()
}

func (r *ScrapedDataRequestPDF) FormatContent() string {
	var formattedContent strings.Builder

	formattedContent.WriteString(fmt.Sprintf("PDF Document: %s\n\n", r.Title))

	if r.Content != nil && *r.Content != "" {
		formattedContent.WriteString("Content:\n")
		formattedContent.WriteString(*r.Content)
	} else {
		formattedContent.WriteString("No content available in this PDF.")
	}

	return formattedContent.String()
}

type Header struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Link struct {
	Text string `json:"text"`
	Href string `json:"href"`
}

type ValidateKeyRequest struct {
	APIKey string `json:"apiKey"`
}
