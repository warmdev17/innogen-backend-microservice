package dto

import "innogen-backend/shared/models"

// RepositoryResponse wraps a repository for HTTP responses.
type RepositoryResponse struct {
	Repository *models.Repository `json:"repository"`
}

// ListRepositoriesResponse wraps a list of repositories.
type ListRepositoriesResponse struct {
	Repositories []models.Repository `json:"repositories"`
}

// ListCommitsResponse wraps a list of commits.
type ListCommitsResponse struct {
	Commits []models.SubmissionCommit `json:"commits"`
}
