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

// GithubConnect handles GET /auth/github/connect
func (h *Handler) GithubConnect(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	email, _ := middleware.GetUserEmail(r)
	url := h.svc.GithubConnectURL(userID, email)
	response.JSON(w, http.StatusOK, dto.GithubConnectResponse{InstallURL: url})
}

// GithubCallback handles GET /auth/github/callback (GitHub redirects here)
func (h *Handler) GithubCallback(w http.ResponseWriter, r *http.Request) {
	installationID := r.URL.Query().Get("installation_id")
	state := r.URL.Query().Get("state")

	if installationID == "" || state == "" {
		http.Redirect(w, r, h.svc.FrontendURL()+"/settings?github_status=error&message=missing_params", http.StatusFound)
		return
	}

	redirectURL, err := h.svc.HandleGithubCallback(r.Context(), installationID, state)
	if err != nil {
		h.log.Error("github callback failed", slog.String("error", err.Error()))
		http.Redirect(w, r, h.svc.FrontendURL()+"/settings?github_status=error&message=server_error", http.StatusFound)
		return
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// Register handles POST /auth/register.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.svc.Register(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			response.Error(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrEmailTaken):
			response.Error(w, http.StatusConflict, "Email already registered")
		case errors.Is(err, service.ErrUsernameTaken):
			response.Error(w, http.StatusConflict, "Username already taken")
		default:
			h.log.Error("register failed", slog.String("error", err.Error()))
			response.Error(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	response.JSON(w, http.StatusCreated, resp)
}

// GithubStatus handles GET /auth/github/status
func (h *Handler) GithubStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	resp, err := h.svc.GithubStatus(r.Context(), userID)
	if err != nil {
		h.log.Error("github status failed", slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	response.JSON(w, http.StatusOK, resp)
}
