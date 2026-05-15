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

// GetProblemByID retrieves a full problem by ID.
// Returns nil, nil if not found.
func (r *SubmissionRepository) GetProblemByID(ctx context.Context, id int) (*models.Problem, error) {
	p := &models.Problem{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, slug, title, difficulty, problem_md, time_limit_ms, memory_limit_mb, acceptance_rate, is_published, sample_test_cases, created_at, updated_at
         FROM problems WHERE id = $1`, id,
	).Scan(&p.ID, &p.Slug, &p.Title, &p.Difficulty, &p.ProblemMD, &p.TimeLimitMs, &p.MemoryLimitMb,
		&p.AcceptanceRate, &p.IsPublished, &p.SampleTestCases, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository.GetProblemByID: %w", err)
	}
	return p, nil
}

// GetLanguageByID retrieves an active language by ID.
// Returns nil, nil if not found.
func (r *SubmissionRepository) GetLanguageByID(ctx context.Context, id int) (*models.Language, error) {
	l := &models.Language{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, piston_alias, piston_version, file_extension, default_file_name, is_active, created_at, updated_at
         FROM languages WHERE id = $1 AND is_active = true`, id,
	).Scan(&l.ID, &l.Name, &l.PistonAlias, &l.PistonVersion, &l.FileExtension, &l.DefaultFileName,
		&l.IsActive, &l.CreatedAt, &l.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository.GetLanguageByID: %w", err)
	}
	return l, nil
}

// GetTestCasesByProblemID retrieves all test cases for a problem, ordered by order_index.
// Returns an empty slice if none found.
func (r *SubmissionRepository) GetTestCasesByProblemID(ctx context.Context, problemID int) ([]models.TestCase, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, problem_id, visibility, input_data, expected_output, execute_code, order_index, created_at
         FROM test_cases WHERE problem_id = $1 ORDER BY order_index ASC`, problemID,
	)
	if err != nil {
		return nil, fmt.Errorf("repository.GetTestCasesByProblemID: %w", err)
	}
	defer rows.Close()

	var testCases []models.TestCase
	for rows.Next() {
		var tc models.TestCase
		if err := rows.Scan(&tc.ID, &tc.ProblemID, &tc.Visibility, &tc.InputData, &tc.ExpectedOutput,
			&tc.ExecuteCode, &tc.OrderIndex, &tc.CreatedAt); err != nil {
			return nil, fmt.Errorf("repository.GetTestCasesByProblemID: %w", err)
		}
		testCases = append(testCases, tc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository.GetTestCasesByProblemID: %w", err)
	}
	return testCases, nil
}

// UpdateSubmissionStatus updates the status of a submission.
func (r *SubmissionRepository) UpdateSubmissionStatus(ctx context.Context, id string, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE submissions SET status = $2 WHERE id = $1`, id, status)
	if err != nil {
		return fmt.Errorf("repository.UpdateSubmissionStatus: %w", err)
	}
	return nil
}

// UpdateSubmissionResult updates the final result of a judged submission.
func (r *SubmissionRepository) UpdateSubmissionResult(ctx context.Context, id string, status string, runtimeMs *int, memoryKb *int, errorMessage *string, passCount int, totalTestcases int) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE submissions
         SET status = $2, runtime_ms = $3, memory_kb = $4, error_message = $5,
             pass_count = $6, total_testcases = $7, judged_at = CURRENT_TIMESTAMP
         WHERE id = $1`,
		id, status, runtimeMs, memoryKb, errorMessage, passCount, totalTestcases,
	)
	if err != nil {
		return fmt.Errorf("repository.UpdateSubmissionResult: %w", err)
	}
	return nil
}
