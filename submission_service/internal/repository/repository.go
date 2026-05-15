package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"innogen-backend/shared/models"
)

// ErrSpamCooldown is returned when the anti-spam trigger prevents a submission.
var ErrSpamCooldown = errors.New("Please wait 10 seconds before submitting again")

// SubmissionRepository handles database queries for the submissions table.
type SubmissionRepository struct {
	pool *pgxpool.Pool
}

// New creates a new SubmissionRepository.
func New(pool *pgxpool.Pool) *SubmissionRepository {
	return &SubmissionRepository{pool: pool}
}

const selectAllSubmissionColumns = `SELECT id, user_id, problem_id, language_id, code, status, runtime_ms, memory_kb, error_message, pass_count, total_testcases, repo_path, commit_sha, created_at, judged_at FROM submissions`

const selectSubmissionColumnsNoCode = `SELECT id, user_id, problem_id, language_id, status, runtime_ms, memory_kb, error_message, pass_count, total_testcases, repo_path, commit_sha, created_at, judged_at FROM submissions`

// CreateSubmission inserts a new Pending submission and returns it.
func (r *SubmissionRepository) CreateSubmission(ctx context.Context, userID, problemID, languageID int, code string) (*models.Submission, error) {
	s := &models.Submission{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO submissions (user_id, problem_id, language_id, code)
         VALUES ($1, $2, $3, $4)
         RETURNING `+selectAllSubmissionColumns,
		userID, problemID, languageID, code,
	).Scan(
		&s.ID, &s.UserID, &s.ProblemID, &s.LanguageID, &s.Code,
		&s.Status, &s.RuntimeMs, &s.MemoryKb, &s.ErrorMessage,
		&s.PassCount, &s.TotalTestcases, &s.RepoPath, &s.CommitSha,
		&s.CreatedAt, &s.JudgedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Please wait 10 seconds") {
			return nil, ErrSpamCooldown
		}
		return nil, fmt.Errorf("repository.CreateSubmission: %w", err)
	}
	return s, nil
}

// FindByID looks up a submission by its UUID.
// Returns nil, nil if not found.
func (r *SubmissionRepository) FindByID(ctx context.Context, id string) (*models.Submission, error) {
	s := &models.Submission{}
	err := r.pool.QueryRow(ctx,
		selectAllSubmissionColumns+" WHERE id = $1",
		id,
	).Scan(
		&s.ID, &s.UserID, &s.ProblemID, &s.LanguageID, &s.Code,
		&s.Status, &s.RuntimeMs, &s.MemoryKb, &s.ErrorMessage,
		&s.PassCount, &s.TotalTestcases, &s.RepoPath, &s.CommitSha,
		&s.CreatedAt, &s.JudgedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository.FindByID: %w", err)
	}
	return s, nil
}

// FindByUserID returns all submissions for a user, newest first, without code.
func (r *SubmissionRepository) FindByUserID(ctx context.Context, userID int) ([]models.Submission, error) {
	rows, err := r.pool.Query(ctx,
		selectSubmissionColumnsNoCode+" WHERE user_id = $1 ORDER BY created_at DESC LIMIT 50",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("repository.FindByUserID: %w", err)
	}
	defer rows.Close()

	var submissions []models.Submission
	for rows.Next() {
		var s models.Submission
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.ProblemID, &s.LanguageID,
			&s.Status, &s.RuntimeMs, &s.MemoryKb, &s.ErrorMessage,
			&s.PassCount, &s.TotalTestcases, &s.RepoPath, &s.CommitSha,
			&s.CreatedAt, &s.JudgedAt,
		); err != nil {
			return nil, fmt.Errorf("repository.FindByUserID: %w", err)
		}
		submissions = append(submissions, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository.FindByUserID: %w", err)
	}
	return submissions, nil
}

// FindLatestByUserAndProblem returns the most recent submission for a user+problem pair.
// Returns nil, nil if none found.
func (r *SubmissionRepository) FindLatestByUserAndProblem(ctx context.Context, userID, problemID int) (*models.Submission, error) {
	s := &models.Submission{}
	err := r.pool.QueryRow(ctx,
		selectAllSubmissionColumns+" WHERE user_id = $1 AND problem_id = $2 ORDER BY created_at DESC LIMIT 1",
		userID, problemID,
	).Scan(
		&s.ID, &s.UserID, &s.ProblemID, &s.LanguageID, &s.Code,
		&s.Status, &s.RuntimeMs, &s.MemoryKb, &s.ErrorMessage,
		&s.PassCount, &s.TotalTestcases, &s.RepoPath, &s.CommitSha,
		&s.CreatedAt, &s.JudgedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository.FindLatestByUserAndProblem: %w", err)
	}
	return s, nil
}

// ProblemExists checks whether a problem with the given ID exists.
func (r *SubmissionRepository) ProblemExists(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM problems WHERE id = $1)`, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("repository.ProblemExists: %w", err)
	}
	return exists, nil
}

// LanguageExists checks whether an active language with the given ID exists.
func (r *SubmissionRepository) LanguageExists(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM languages WHERE id = $1 AND is_active = true)`, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("repository.LanguageExists: %w", err)
	}
	return exists, nil
}
