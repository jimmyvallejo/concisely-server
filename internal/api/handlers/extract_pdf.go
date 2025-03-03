package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func (h *Handlers) ExtractPDF(w http.ResponseWriter, r *http.Request) {
	request := extractPDFRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	resp, err := http.Get(request.URL)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Resource not available for download")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Failed to download PDF, status code: %d", resp.StatusCode))
		return
	}

	tempDir, err := os.MkdirTemp("", "pdf-extract-*")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create temporary directory")
		return
	}
	defer os.RemoveAll(tempDir) 

	tempFile, err := os.CreateTemp(tempDir, "downloaded-*.pdf")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Temporary storage not operational")
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to copy PDF contents")
		return
	}

	tempFile.Sync()

	extractDir := filepath.Join(tempDir, "extracted")
	err = os.Mkdir(extractDir, 0755)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create extraction directory")
		return
	}

	conf := api.LoadConfiguration()
	err = api.ExtractContentFile(tempFile.Name(), extractDir, nil, conf)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to extract content: %v", err))
		return
	}

	var combinedContent strings.Builder
	files, err := os.ReadDir(extractDir)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to read extracted content")
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			content, err := os.ReadFile(filepath.Join(extractDir, file.Name()))
			if err != nil {
				continue
			}
			combinedContent.WriteString(string(content))
			combinedContent.WriteString("\n")
		}
	}

	response := ExtractedPDFResponse{
		Content: combinedContent.String(),
	}

	respondWithJSON(w, http.StatusOK, response)

}
