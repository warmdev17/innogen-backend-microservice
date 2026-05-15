package dto

import (
	"strings"
	"time"

	"innogen-backend/shared/models"
)

// SubmitRequest is the JSON body for POST /submit.
type SubmitRequest struct {
	ProblemID  int    `json:"problemId"`
	LanguageID int    `json:"languageId"`
	Code       string `json:"code"`
}

// Validate returns an error message if invalid, or empty string if valid.
func (r *SubmitRequest) Validate() string {
	if r.ProblemID <= 0 {
		return "Problem ID is required"
	}
	if r.LanguageID <= 0 {
		return "Language ID is required"
	}
	if strings.TrimSpace(r.Code) == "" {
		return "Code must not be empty"
	}
	return ""
}

// SubmitResponse wraps a submission for POST /submit response.
type SubmitResponse struct {
	Submission *models.Submission `json:"submission"`
}

// GetSubmissionResponse wraps a submission for GET /submissions/{id} response.
type GetSubmissionResponse struct {
	Submission *models.Submission `json:"submission"`
}

// SubmissionListItem is a submission without the code field, for list endpoints.
type SubmissionListItem struct {
	ID             string  `json:"id"`
	UserID         int     `json:"userId"`
	ProblemID      int     `json:"problemId"`
	LanguageID     int     `json:"languageId"`
	Status         string  `json:"status"`
	RuntimeMs      *int    `json:"runtimeMs"`
	MemoryKb       *int    `json:"memoryKb"`
	ErrorMessage   *string `json:"errorMessage"`
	PassCount      int     `json:"passCount"`
	TotalTestcases int     `json:"totalTestcases"`
	RepoPath       *string `json:"repoPath"`
	CommitSha      *string `json:"commitSha"`
	CreatedAt      string  `json:"createdAt"`
	JudgedAt       *string `json:"judgedAt"`
}

// ToSubmissionListItem converts a models.Submission to a SubmissionListItem.
func ToSubmissionListItem(s *models.Submission) SubmissionListItem {
	item := SubmissionListItem{
		ID:             s.ID,
		UserID:         s.UserID,
		ProblemID:      s.ProblemID,
		LanguageID:     s.LanguageID,
		Status:         s.Status,
		RuntimeMs:      s.RuntimeMs,
		MemoryKb:       s.MemoryKb,
		ErrorMessage:   s.ErrorMessage,
		PassCount:      s.PassCount,
		TotalTestcases: s.TotalTestcases,
		RepoPath:       s.RepoPath,
		CommitSha:      s.CommitSha,
		CreatedAt:      s.CreatedAt.Format(time.RFC3339),
	}
	if s.JudgedAt != nil {
		js := s.JudgedAt.Format(time.RFC3339)
		item.JudgedAt = &js
	}
	return item
}

// ListSubmissionsResponse wraps submission list for GET /me/submissions.
type ListSubmissionsResponse struct {
	Submissions []SubmissionListItem `json:"submissions"`
}

// LatestSubmissionResponse wraps a submission for GET /me/submissions/{problemId}/latest.
type LatestSubmissionResponse struct {
	Submission *models.Submission `json:"submission"`
}
