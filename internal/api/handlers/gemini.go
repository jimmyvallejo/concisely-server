package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

		writeSSE(w, flusher, "Connection established. Processing started...")

		var errorSent bool
		chunkCount := 0

		for {
			select {
			case <-ctx.Done():
				if !errorSent {
					writeSSE(w, flusher, "Request timed out or canceled")
					log.Println("Client connection closed or timed out")
				}
				return

			case err := <-errorChan:
				errorSent = true
				log.Printf("Error: %v", err)
				writeSSE(w, flusher, fmt.Sprintf("Error: %v", err))
				time.Sleep(100 * time.Millisecond)
				return

			case <-doneChan:
				writeSSE(w, flusher, "[DONE]")
				log.Printf("PDF processing complete. Sent %d chunks in total.", chunkCount)
				return

			case msg := <-messageChan:
				chunkCount++
				log.Printf("CHUNK #%d: %s", chunkCount, msg)
				
				writeSSE(w, flusher, msg)

			case <-time.After(30 * time.Second):
				writeSSE(w, flusher, "still processing...")
			}
		}
	}()

	<-ctx.Done()
}

func writeSSE(w http.ResponseWriter, flusher http.Flusher, message string) {
	if message == "" {
		return
	}

	_, err := fmt.Fprintf(w, "data: %s\n\n", message)
	if err != nil {
		log.Printf("Error writing to response: %v", err)
		return
	}

	flusher.Flush()
}

func processPDF(ctx context.Context, request GeminiPDFRequest, messageChan chan<- string, doneChan chan<- bool, errorChan chan<- error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in processPDF: %v", r)
			errorChan <- fmt.Errorf("internal server error: %v", r)
		}
	}()

	client, err := genai.NewClient(ctx, option.WithAPIKey(request.ApiKey))
	if err != nil {
		log.Printf("Failed to create Gemini client: %v", err)
		errorChan <- fmt.Errorf("invalid API key: %v", err)
		return
	}
	defer client.Close()

	messageChan <- "API connection established. Downloading PDF..."

	pdfBytes, err := downloadPDF(ctx, request.URL)
	if err != nil {
		errorChan <- fmt.Errorf("failed to download PDF: %v", err)
		return
	}

	messageChan <- fmt.Sprintf("PDF downloaded (%d bytes). Processing with Gemini...", len(pdfBytes))

	model := client.GenerativeModel("gemini-2.0-flash")
	req := []genai.Part{
		genai.Blob{MIMEType: "application/pdf", Data: pdfBytes},
		genai.Text(systemPromptPDF),
	}

	iter := model.GenerateContentStream(ctx, req...)

	for {
		select {
		case <-ctx.Done():
			errorChan <- fmt.Errorf("context canceled or timed out")
			return

		default:
			resp, err := iter.Next()
			if err == iterator.Done {
				doneChan <- true
				return
			}
			if err != nil {
				errorChan <- fmt.Errorf("error generating content: %v", err)
				return
			}

			for _, c := range resp.Candidates {
				if c.Content != nil {
					for _, part := range c.Content.Parts {
						partStr := fmt.Sprintf("%v", part)
						messageChan <- partStr
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