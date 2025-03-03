package routes

import (
	"net/http"

	"github.com/jimmyvallejo/concisely-server/internal/api/handlers"
)

type Router struct {
	mux      *http.ServeMux
	handlers *handlers.Handlers
}

func NewRouter(h *handlers.Handlers) *Router {
	return &Router{
		mux:      http.NewServeMux(),
		handlers: h,
	}
}

func (r *Router) SetupRoutes() {

	// System Health
	r.mux.HandleFunc("GET /healthz", handlers.HandlerReadiness)

	// Anthropic
	r.mux.HandleFunc("POST /anthropic-stream", r.handlers.AnthropicCompletion)
	r.mux.HandleFunc("POST /anthropic-validate", r.handlers.ValidateAnthropicKey)

	// ChatGPT
	r.mux.HandleFunc("POST /gpt-stream", r.handlers.ChatGPTCompletion)
	r.mux.HandleFunc("POST /gpt-validate", r.handlers.ValidateOpenAIKey)

	// Deepseek
	r.mux.HandleFunc("POST /deepseek-validate", r.handlers.ValidateDeepseekKey)

	//PDF Extraction
	r.mux.HandleFunc("POST /extract-pdf", r.handlers.ExtractPDF)

}

func (r *Router) GetHandler() http.Handler {
	return r.mux
}
