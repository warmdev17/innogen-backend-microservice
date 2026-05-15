package route

import (
	"net/http"

	"innogen-backend/run_service/internal/handler"
)

// Register registers all run service routes on the given ServeMux.
func Register(mux *http.ServeMux, h *handler.Handler) {
	mux.HandleFunc("POST /run", h.Run)
}
