package problem

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"innogen-backend/shared/models"
)

// ProblemRepository provides database access for problem entities.
type ProblemRepository struct {
	pool *pgxpool.Pool
}

// NewProblemRepository creates a new ProblemRepository.
func NewProblemRepository(pool *pgxpool.Pool) *ProblemRepository {
	return &ProblemRepository{pool: pool}
}

// FindBySlug returns a single published problem by slug.
// Returns nil, nil if not found.
func (r *ProblemRepository) FindBySlug(ctx context.Context, slug string) (*models.Problem, error) {
	query := `
		SELECT id, slug, title, difficulty, problem_md, time_limit_ms, memory_limit_mb,
		       acceptance_rate, is_published, execution_mode, function_name, initial_code, driver_code, solution_file_name, sample_test_cases, created_at, updated_at
		FROM problems
		WHERE slug = $1 AND is_published = true
	`

	var p models.Problem
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&p.ID, &p.Slug, &p.Title, &p.Difficulty, &p.ProblemMD,
		&p.TimeLimitMs, &p.MemoryLimitMb, &p.AcceptanceRate,
		&p.IsPublished, &p.ExecutionMode, &p.FunctionName, &p.InitialCode, &p.DriverCode, &p.SolutionFileName, &p.SampleTestCases, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find problem by slug %q: %w", slug, err)
	}

	return &p, nil
}

// FindByID returns a single published problem by id.
// Returns nil, nil if not found.
func (r *ProblemRepository) FindByID(ctx context.Context, id int) (*models.Problem, error) {
	query := `
		SELECT id, slug, title, difficulty, problem_md, time_limit_ms, memory_limit_mb,
		       acceptance_rate, is_published, execution_mode, function_name, initial_code, driver_code, solution_file_name, sample_test_cases, created_at, updated_at
		FROM problems
		WHERE id = $1 AND is_published = true
	`

	var p models.Problem
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Slug, &p.Title, &p.Difficulty, &p.ProblemMD,
		&p.TimeLimitMs, &p.MemoryLimitMb, &p.AcceptanceRate,
		&p.IsPublished, &p.ExecutionMode, &p.FunctionName, &p.InitialCode, &p.DriverCode, &p.SolutionFileName, &p.SampleTestCases, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("find problem by id %d: %w", id, err)
	}

	return &p, nil
}

// FindTestCasesByProblemID returns test cases for a given problem filtered by visibility.
// For public API usage, visibility should be "sample".
// Hidden test cases are never returned unless visibility="hidden" is explicitly passed.
func (r *ProblemRepository) FindTestCasesByProblemID(ctx context.Context, problemID int, visibility string) ([]models.TestCase, error) {
	query := `
		SELECT id, problem_id, visibility, input_data, expected_output, order_index
		FROM test_cases
		WHERE problem_id = $1 AND visibility = $2
		ORDER BY order_index ASC
	`

	rows, err := r.pool.Query(ctx, query, problemID, visibility)
	if err != nil {
		return nil, fmt.Errorf("find test cases by problem id %d: %w", problemID, err)
	}
	defer rows.Close()

	var cases []models.TestCase
	for rows.Next() {
		var tc models.TestCase
		if err := rows.Scan(
			&tc.ID, &tc.ProblemID, &tc.Visibility, &tc.InputData, &tc.ExpectedOutput, &tc.OrderIndex,
		); err != nil {
			return nil, fmt.Errorf("scan test case row: %w", err)
		}
		cases = append(cases, tc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate test case rows: %w", err)
	}

	return cases, nil
}
