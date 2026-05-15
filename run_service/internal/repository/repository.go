package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"innogen-backend/shared/models"
)

// Repository handles database queries for the run service.
type Repository struct {
	pool *pgxpool.Pool
}

// New creates a new Repository with the given connection pool.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// ProblemExists checks whether a problem with the given ID exists.
func (r *Repository) ProblemExists(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM problems WHERE id = $1)`,
		id,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("repository.ProblemExists: %w", err)
	}
	return exists, nil
}

// GetLanguageByID looks up an active language by its ID.
// Returns nil, nil if no active language is found.
func (r *Repository) GetLanguageByID(ctx context.Context, id int) (*models.Language, error) {
	l := &models.Language{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, piston_alias, piston_version, file_extension, default_file_name
		 FROM languages WHERE id = $1 AND is_active = true`,
		id,
	).Scan(&l.ID, &l.Name, &l.PistonAlias, &l.PistonVersion, &l.FileExtension, &l.DefaultFileName)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository.GetLanguageByID: %w", err)
	}
	return l, nil
}

// GetSampleTestCases retrieves all sample test cases for a problem, ordered by order_index.
// Returns an empty slice if no test cases are found.
func (r *Repository) GetSampleTestCases(ctx context.Context, problemID int) ([]models.TestCase, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, problem_id, visibility, input_data, expected_output, order_index
		 FROM test_cases WHERE problem_id = $1 AND visibility = 'sample'
		 ORDER BY order_index ASC`,
		problemID,
	)
	if err != nil {
		return nil, fmt.Errorf("repository.GetSampleTestCases: %w", err)
	}
	defer rows.Close()

	var testCases []models.TestCase
	for rows.Next() {
		var tc models.TestCase
		if err := rows.Scan(&tc.ID, &tc.ProblemID, &tc.Visibility, &tc.InputData, &tc.ExpectedOutput, &tc.OrderIndex); err != nil {
			return nil, fmt.Errorf("repository.GetSampleTestCases: %w", err)
		}
		testCases = append(testCases, tc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository.GetSampleTestCases: %w", err)
	}
	return testCases, nil
}
