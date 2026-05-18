package route

import (
	"net/http"

	"innogen-backend/repo_service/internal/handler"
	"innogen-backend/repo_service/internal/webhook"
	"innogen-backend/shared/middleware"
)

// Register registers all repo service routes on the given ServeMux.
func Register(mux *http.ServeMux, h *handler.Handler, wh *webhook.WebhookHandler) {
	mux.Handle("GET /repositories", middleware.XUserID()(http.HandlerFunc(h.ListRepositories)))
	mux.Handle("GET /repositories/{id}/commits", middleware.XUserID()(http.HandlerFunc(h.ListCommits)))
	mux.HandleFunc("POST /internal/commits/accepted-submission", h.CommitAcceptedSubmission)
	mux.HandleFunc("POST /webhooks/github", wh.HandleWebhook)
}
