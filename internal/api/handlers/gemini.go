package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func (h *Handlers) GeminiParsePDF(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	flusher, ok := w.(http.Flusher)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "streaming not supported")
		return
	}

	if err := h.Limiter.Wait(ctx); err != nil {
		respondWithError(w, http.StatusTooManyRequests, "Rate limit exceeded")
		return
	}

	var request GeminiPDFRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to read request body: "+err.Error())
		return
	}

	if request.ApiKey == "" || request.URL == "" {
		respondWithError(w, http.StatusBadRequest, "API key and URL are required")
		return
	}

	messageChan := make(chan string, 100)
	doneChan := make(chan bool, 1)
	errorChan := make(chan error, 1)

	go processPDF(ctx, request, messageChan, doneChan, errorChan)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic in writer goroutine: %v", r)
			}
		}()

		var errorSent bool
		chunkCount := 0

		for {
			select {
			case <-ctx.Done():
				if !errorSent {
					if !isConnectionClosed(r) {
						writeSSE(w, flusher, "Request timed out or canceled")
					}
					log.Println("Client connection closed or timed out")
				}
				return

			case err := <-errorChan:
				if err == nil {
					continue 
				}
				errorSent = true
				log.Printf("Error: %v", err)
				if !isConnectionClosed(r) {
					writeSSE(w, flusher, fmt.Sprintf("Error: %v", err))
				}
				time.Sleep(100 * time.Millisecond)
				return

			case <-doneChan:
				return

			case msg := <-messageChan:
				if msg == "" {
					continue 
				}
				chunkCount++
				if !isConnectionClosed(r) {
					if !writeSSE(w, flusher, msg) {
						log.Println("Failed to write to client, connection likely closed")
						return
					}
				} else {
					log.Println("Connection closed, stopping message processing")
					return
				}

			case <-time.After(30 * time.Second):
				if !isConnectionClosed(r) {
					if !writeSSE(w, flusher, "still processing...") {
						log.Println("Failed to write keep-alive, connection likely closed")
						return
					}
				} else {
					log.Println("Connection closed during keep-alive")
					return
				}
			}
		}
	}()

	<-ctx.Done()
}

func writeSSE(w http.ResponseWriter, flusher http.Flusher, message string) bool {
	if message == "" {
		return true
	}

	_, err := fmt.Fprintf(w, message)
	if err != nil {
		log.Printf("Error writing to response: %v", err)
		return false
	}

	flusher.Flush()
	return true
}

func isConnectionClosed(r *http.Request) bool {
	select {
	case <-r.Context().Done():
		return true
	default:
		return false
	}
}

func processPDF(ctx context.Context, request GeminiPDFRequest, messageChan chan<- string, doneChan chan<- bool, errorChan chan<- error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in processPDF: %v", r)
			select {
			case errorChan <- fmt.Errorf("internal server error: %v", r):
			default:
			}
		}
	}()

	client, err := genai.NewClient(ctx, option.WithAPIKey(request.ApiKey))
	if err != nil {
		log.Printf("Failed to create Gemini client: %v", err)
		select {
		case errorChan <- fmt.Errorf("invalid API key: %v", err):
		case <-ctx.Done():
		}
		return
	}
	defer client.Close()

	log.Println("Downloading PDF...")

	pdfBytes, err := downloadPDF(ctx, request.URL)
	if err != nil {
		select {
		case errorChan <- fmt.Errorf("failed to download PDF: %v", err):
		case <-ctx.Done():
		}
		return
	}

	model := client.GenerativeModel("gemini-2.0-flash")
	req := []genai.Part{
		genai.Blob{MIMEType: "application/pdf", Data: pdfBytes},
		genai.Text(systemPromptPDF),
	}

	iter := model.GenerateContentStream(ctx, req...)

	var contentBuffer strings.Builder

	for {
		select {
		case <-ctx.Done():
			select {
			case errorChan <- fmt.Errorf("context canceled or timed out"):
			default:
			}
			return

		default:
			resp, err := iter.Next()
			if err == iterator.Done {
				if contentBuffer.Len() > 0 {
					select {
					case messageChan <- contentBuffer.String():
					case <-ctx.Done():
						return
					}
				}
				select {
				case doneChan <- true:
				case <-ctx.Done():
				}
				return
			}
			if err != nil {
				select {
				case errorChan <- fmt.Errorf("error generating content: %v", err):
				case <-ctx.Done():
				}
				return
			}
			for _, c := range resp.Candidates {
				if c.Content != nil {
					for _, part := range c.Content.Parts {
						partStr := fmt.Sprintf("%v", part)
						select {
						case messageChan <- partStr:
							time.Sleep(100 * time.Millisecond)
							contentBuffer.WriteString(partStr)
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}
	}
}

func downloadPDF(ctx context.Context, url string) ([]byte, error) {

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	pdfReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	pdfResp, err := httpClient.Do(pdfReq)
	if err != nil {
		return nil, fmt.Errorf("failed to download: %v", err)
	}
	defer pdfResp.Body.Close()

	if pdfResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status: %d", pdfResp.StatusCode)
	}

	const maxPDFSize = 50 * 1024 * 1024
	return io.ReadAll(io.LimitReader(pdfResp.Body, maxPDFSize))
}

func (h *Handlers) ValidateGeminiKey(w http.ResponseWriter, r *http.Request) {
	request := ValidateKeyRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	client, err := genai.NewClient(context.Background(), option.WithAPIKey(request.APIKey))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create Gemini client")
		return
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash")
	model.SetMaxOutputTokens(1)

	_, err = model.GenerateContent(context.Background(), genai.Text("What is a quaternion?"))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	respondNoBody(w, http.StatusOK)
}
