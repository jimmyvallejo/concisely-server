package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (h *Handlers) ValidateDeepseekKey(w http.ResponseWriter, r *http.Request) {
	request := ValidateKeyRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", strings.NewReader(`{
        "model": "deepseek-chat",
        "messages": [{"role": "user", "content": "test"}],
        "max_tokens": 1
    }`))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create request")
		return
	}

	req.Header.Add("Authorization", "Bearer "+request.APIKey)
	req.Header.Add("Content-Type", "application/json")

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
