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

	// ChatGPT
	r.mux.HandleFunc("GET /gpt-test", r.handlers.ChatGPTSample)

}

func (r *Router) GetHandler() http.Handler {
	return r.mux
}
