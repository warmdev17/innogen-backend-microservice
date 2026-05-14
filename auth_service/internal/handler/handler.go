package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"innogen-backend/auth_service/internal/dto"
	"innogen-backend/auth_service/internal/service"
	"innogen-backend/shared/middleware"
	"innogen-backend/shared/response"
)

// Handler holds dependencies for HTTP request handlers.
type Handler struct {
	svc *service.AuthService
	log *slog.Logger
}

// New creates a new Handler with the given dependencies.
func New(svc *service.AuthService, log *slog.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// Login handles POST /auth/login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.svc.Login(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			response.Error(w, http.StatusBadRequest, "Invalid request body")
		case errors.Is(err, service.ErrInvalidCredentials):
			response.Error(w, http.StatusUnauthorized, "Invalid email or password")
		case errors.Is(err, service.ErrAccountInactive):
			response.Error(w, http.StatusForbidden, "Account is inactive")
		default:
			h.log.Error("login failed", slog.String("error", err.Error()))
			response.Error(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	response.JSON(w, http.StatusOK, resp)
}

// CurrentUser handles GET /auth/me.
func (h *Handler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	resp, err := h.svc.CurrentUser(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			response.Error(w, http.StatusUnauthorized, "User not found")
		case errors.Is(err, service.ErrAccountInactive):
			response.Error(w, http.StatusForbidden, "Account is inactive")
		default:
			h.log.Error("current user lookup failed", slog.String("error", err.Error()))
			response.Error(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	response.JSON(w, http.StatusOK, resp)
}
