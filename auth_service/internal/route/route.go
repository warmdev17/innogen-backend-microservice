package route

import (
	"net/http"

	"innogen-backend/auth_service/internal/handler"
	"innogen-backend/shared/middleware"
)

// Register registers all auth service routes on the given ServeMux.
func Register(mux *http.ServeMux, h *handler.Handler, jwtSecret string) {
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.Handle("GET /auth/me", middleware.Auth(jwtSecret)(http.HandlerFunc(h.CurrentUser)))
}
