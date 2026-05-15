package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"innogen-backend/repo_service/internal/dto"
	"innogen-backend/repo_service/service"
	"innogen-backend/shared/middleware"
	"innogen-backend/shared/response"
)

// Handler holds dependencies for HTTP request handlers.
type Handler struct {
	svc *service.RepoService
	log *slog.Logger
}

// New creates a new Handler.
func New(svc *service.RepoService, log *slog.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// ListRepositories handles GET /repositories.
func (h *Handler) ListRepositories(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	repos, err := h.svc.ListRepositories(r.Context(), userID)
	if err != nil {
		h.log.Error("list repositories failed", slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response.JSON(w, http.StatusOK, dto.ListRepositoriesResponse{Repositories: repos})
}

// CommitAcceptedSubmission handles POST /internal/commits/accepted-submission.
func (h *Handler) CommitAcceptedSubmission(w http.ResponseWriter, r *http.Request) {
	var req dto.CommitAcceptedSubmissionRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.SubmissionID == "" {
		response.Error(w, http.StatusBadRequest, "Missing required field: submissionId")
		return
	}
	if req.UserID <= 0 {
		response.Error(w, http.StatusBadRequest, "Missing required field: userId")
		return
	}
	if req.ProblemID <= 0 {
		response.Error(w, http.StatusBadRequest, "Missing required field: problemId")
		return
	}
	if req.LanguageID <= 0 {
		response.Error(w, http.StatusBadRequest, "Missing required field: languageId")
		return
	}
	if req.Code == "" {
		response.Error(w, http.StatusBadRequest, "Missing required field: code")
		return
	}

	svcReq := service.CommitAcceptedSubmissionRequest{
		SubmissionID: req.SubmissionID,
		UserID:       req.UserID,
		ProblemID:    req.ProblemID,
		LanguageID:   req.LanguageID,
		Code:         req.Code,
	}

	resp, err := h.svc.CommitAcceptedSubmission(r.Context(), svcReq)
	if err != nil {
		if errors.Is(err, service.ErrGithubAccountNotFound) {
			response.Error(w, http.StatusNotFound, "GitHub account not linked")
			return
		}
		h.log.Error("commit accepted submission failed", slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response.JSON(w, http.StatusOK, resp)
}

// ListCommits handles GET /repositories/{id}/commits.
func (h *Handler) ListCommits(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := r.PathValue("id")
	repoID, err := strconv.Atoi(idStr)
	if err != nil || repoID <= 0 {
		response.Error(w, http.StatusBadRequest, "Invalid repository ID")
		return
	}

	commits, err := h.svc.ListCommits(r.Context(), userID, repoID)
	if err != nil {
		if errors.Is(err, service.ErrRepositoryNotFound) {
			response.Error(w, http.StatusNotFound, "Repository not found")
			return
		}
		h.log.Error("list commits failed", slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response.JSON(w, http.StatusOK, dto.ListCommitsResponse{Commits: commits})
}
