package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"innogen-backend/repo_service/internal/pathbuilder"
	"innogen-backend/repo_service/repository"
	"innogen-backend/shared/languageutil"
	"innogen-backend/shared/models"
)

var ErrRepositoryNotFound = errors.New("repository not found")

// RepoService handles repository and commit business logic.
type RepoService struct {
	repo *repository.RepoRepository
	log  *slog.Logger
}

// New creates a new RepoService.
func New(repo *repository.RepoRepository, log *slog.Logger) *RepoService {
	return &RepoService{repo: repo, log: log}
}

// CommitSubmission triggers the mock commit flow for an Accepted submission.
// This is a best-effort operation: errors are logged but should not fail the caller's pipeline.
func (s *RepoService) CommitSubmission(ctx context.Context, submissionID string, userID, problemID, languageID int) error {
	// Idempotency guard: if already committed, skip
	_, commitSHA, err := s.repo.GetSubmissionRepoInfo(ctx, submissionID)
	if err != nil {
		return err
	}
	if commitSHA != "" {
		s.log.Info("submission already committed, skipping",
			slog.String("submissionId", submissionID),
		)
		return nil
	}

	// 1. Load curriculum context
	curriculum, err := s.repo.GetCurriculumContext(ctx, problemID)
	if err != nil {
		return err
	}

	// 2. Load language info
	lang, err := s.repo.GetLanguageByID(ctx, languageID)
	if err != nil {
		return err
	}
	if lang == nil {
		return fmt.Errorf("language %d not found", languageID)
	}

	// 3. Build paths
	fileName := languageutil.DetermineFileName(lang)
	newRepoPath := pathbuilder.BuildFilePath(*curriculum, fileName)
	repoName := pathbuilder.BuildRepoName(curriculum.SubjectSlug)
	newCommitSHA := pathbuilder.GenerateCommitSHA()

	// 4-6. Execute commit operations in a transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	repoID, err := s.repo.UpsertRepositoryTx(ctx, tx, userID, curriculum.SubjectID, repoName)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateSubmissionRepoInfoTx(ctx, tx, submissionID, newRepoPath, newCommitSHA); err != nil {
		return err
	}
	if err := s.repo.InsertSubmissionCommitTx(ctx, tx, submissionID, repoID, newRepoPath, newCommitSHA); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	s.log.Info("mock commit created",
		slog.String("submissionId", submissionID),
		slog.String("repoPath", newRepoPath),
		slog.String("commitSha", newCommitSHA),
		slog.Int("repoId", repoID),
	)

	return nil
}

// ListRepositories returns all repositories for a user.
func (s *RepoService) ListRepositories(ctx context.Context, userID int) ([]models.Repository, error) {
	return s.repo.FindRepositoriesByUserID(ctx, userID)
}

// ListCommits returns all commits for a repository, after verifying ownership.
func (s *RepoService) ListCommits(ctx context.Context, userID, repositoryID int) ([]models.SubmissionCommit, error) {
	repo, err := s.repo.FindRepositoryByID(ctx, repositoryID)
	if err != nil {
		return nil, err
	}
	if repo == nil || repo.UserID != userID {
		return nil, ErrRepositoryNotFound
	}
	return s.repo.FindCommitsByRepositoryID(ctx, repositoryID)
}
