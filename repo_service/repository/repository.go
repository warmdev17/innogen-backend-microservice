package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"innogen-backend/repo_service/internal/pathbuilder"
	"innogen-backend/shared/models"
)

// RepoRepository handles database queries for repository and commit operations.
type RepoRepository struct {
	pool *pgxpool.Pool
}

// New creates a new RepoRepository.
func New(pool *pgxpool.Pool) *RepoRepository {
	return &RepoRepository{pool: pool}
}

// GetCurriculumContext loads the curriculum hierarchy for a problem.
// Returns nil, error if the problem is not linked to any lesson.
func (r *RepoRepository) GetCurriculumContext(ctx context.Context, problemID int) (*pathbuilder.CurriculumContext, error) {
	c := &pathbuilder.CurriculumContext{}
	err := r.pool.QueryRow(ctx,
		`SELECT subj.slug, subj.id, ss.order_index, l.order_index, lp.order_index, p.slug
		 FROM lesson_problems lp
		 JOIN lessons l ON lp.lesson_id = l.id
		 JOIN subject_sessions ss ON l.subject_session_id = ss.id
		 JOIN subjects subj ON ss.subject_id = subj.id
		 JOIN problems p ON lp.problem_id = p.id
		 WHERE lp.problem_id = $1
		 ORDER BY lp.order_index ASC
		 LIMIT 1`,
		problemID,
	).Scan(&c.SubjectSlug, &c.SubjectID, &c.SessionOrder, &c.LessonOrder, &c.ProblemOrder, &c.ProblemSlug)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("problem %d is not linked to any lesson", problemID)
	}
	if err != nil {
		return nil, fmt.Errorf("repository.GetCurriculumContext: %w", err)
	}
	return c, nil
}

// GetLanguageByID retrieves an active language by ID.
// Returns nil, nil if not found.
func (r *RepoRepository) GetLanguageByID(ctx context.Context, id int) (*models.Language, error) {
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

// UpsertRepository creates or updates a repository for a user+subject pair, returning the repository ID.
func (r *RepoRepository) UpsertRepository(ctx context.Context, userID, subjectID int, repoName string) (int, error) {
	var repoID int
	err := r.pool.QueryRow(ctx,
		`INSERT INTO repositories (user_id, subject_id, repo_name)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, subject_id) DO UPDATE SET updated_at = CURRENT_TIMESTAMP
		 RETURNING id`,
		userID, subjectID, repoName,
	).Scan(&repoID)
	if err != nil {
		return 0, fmt.Errorf("repository.UpsertRepository: %w", err)
	}
	return repoID, nil
}

// UpdateSubmissionRepoInfo updates the repo_path and commit_sha on a submission.
func (r *RepoRepository) UpdateSubmissionRepoInfo(ctx context.Context, submissionID, repoPath, commitSHA string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE submissions SET repo_path = $2, commit_sha = $3 WHERE id = $1`,
		submissionID, repoPath, commitSHA,
	)
	if err != nil {
		return fmt.Errorf("repository.UpdateSubmissionRepoInfo: %w", err)
	}
	return nil
}

// InsertSubmissionCommit inserts a row into the submission_commits table.
func (r *RepoRepository) InsertSubmissionCommit(ctx context.Context, submissionID string, repositoryID int, filePath, commitSHA string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO submission_commits (submission_id, repository_id, file_path, commit_sha)
		 VALUES ($1, $2, $3, $4)`,
		submissionID, repositoryID, filePath, commitSHA,
	)
	if err != nil {
		return fmt.Errorf("repository.InsertSubmissionCommit: %w", err)
	}
	return nil
}

// FindRepositoryByID looks up a repository by its ID.
// Returns nil, nil if not found.
func (r *RepoRepository) FindRepositoryByID(ctx context.Context, id int) (*models.Repository, error) {
	repo := &models.Repository{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, subject_id, repo_name, repo_full_name, repo_url, github_repo_id, default_branch, created_at, updated_at
		 FROM repositories WHERE id = $1`, id,
	).Scan(&repo.ID, &repo.UserID, &repo.SubjectID, &repo.RepoName,
		&repo.RepoFullName, &repo.RepoURL, &repo.GithubRepoID, &repo.DefaultBranch,
		&repo.CreatedAt, &repo.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository.FindRepositoryByID: %w", err)
	}
	return repo, nil
}

// GetSubmissionRepoInfo returns the repo_path and commit_sha for a submission.
// Returns empty strings if not set.
func (r *RepoRepository) GetSubmissionRepoInfo(ctx context.Context, submissionID string) (string, string, error) {
	var repoPath, commitSHA *string
	err := r.pool.QueryRow(ctx,
		`SELECT repo_path, commit_sha FROM submissions WHERE id = $1`, submissionID,
	).Scan(&repoPath, &commitSHA)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", nil
	}
	if err != nil {
		return "", "", fmt.Errorf("repository.GetSubmissionRepoInfo: %w", err)
	}
	rp := ""
	cs := ""
	if repoPath != nil {
		rp = *repoPath
	}
	if commitSHA != nil {
		cs = *commitSHA
	}
	return rp, cs, nil
}

// BeginTx starts a new database transaction.
func (r *RepoRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

// UpsertRepositoryTx is the transactional version of UpsertRepository.
func (r *RepoRepository) UpsertRepositoryTx(ctx context.Context, tx pgx.Tx, userID, subjectID int, repoName string) (int, error) {
	var repoID int
	err := tx.QueryRow(ctx,
		`INSERT INTO repositories (user_id, subject_id, repo_name)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, subject_id) DO UPDATE SET updated_at = CURRENT_TIMESTAMP
		 RETURNING id`,
		userID, subjectID, repoName,
	).Scan(&repoID)
	if err != nil {
		return 0, fmt.Errorf("repository.UpsertRepositoryTx: %w", err)
	}
	return repoID, nil
}

// UpdateSubmissionRepoInfoTx is the transactional version of UpdateSubmissionRepoInfo.
func (r *RepoRepository) UpdateSubmissionRepoInfoTx(ctx context.Context, tx pgx.Tx, submissionID, repoPath, commitSHA string) error {
	_, err := tx.Exec(ctx,
		`UPDATE submissions SET repo_path = $2, commit_sha = $3 WHERE id = $1`,
		submissionID, repoPath, commitSHA,
	)
	if err != nil {
		return fmt.Errorf("repository.UpdateSubmissionRepoInfoTx: %w", err)
	}
	return nil
}

// InsertSubmissionCommitTx is the transactional version of InsertSubmissionCommit.
func (r *RepoRepository) InsertSubmissionCommitTx(ctx context.Context, tx pgx.Tx, submissionID string, repositoryID int, filePath, commitSHA string) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO submission_commits (submission_id, repository_id, file_path, commit_sha)
		 VALUES ($1, $2, $3, $4)`,
		submissionID, repositoryID, filePath, commitSHA,
	)
	if err != nil {
		return fmt.Errorf("repository.InsertSubmissionCommitTx: %w", err)
	}
	return nil
}

// FindRepositoriesByUserID returns all repositories for a user.
func (r *RepoRepository) FindRepositoriesByUserID(ctx context.Context, userID int) ([]models.Repository, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, subject_id, repo_name, repo_full_name, repo_url, github_repo_id, default_branch, created_at, updated_at
		 FROM repositories WHERE user_id = $1 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("repository.FindRepositoriesByUserID: %w", err)
	}
	defer rows.Close()

	var repos []models.Repository
	for rows.Next() {
		var repo models.Repository
		if err := rows.Scan(&repo.ID, &repo.UserID, &repo.SubjectID, &repo.RepoName,
			&repo.RepoFullName, &repo.RepoURL, &repo.GithubRepoID, &repo.DefaultBranch,
			&repo.CreatedAt, &repo.UpdatedAt); err != nil {
			return nil, fmt.Errorf("repository.FindRepositoriesByUserID: %w", err)
		}
		repos = append(repos, repo)
	}
	return repos, rows.Err()
}

// FindCommitsByRepositoryID returns all commits for a repository.
func (r *RepoRepository) FindCommitsByRepositoryID(ctx context.Context, repositoryID int) ([]models.SubmissionCommit, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, submission_id, repository_id, file_path, commit_sha, commit_url, created_at
		 FROM submission_commits WHERE repository_id = $1 ORDER BY created_at DESC`, repositoryID,
	)
	if err != nil {
		return nil, fmt.Errorf("repository.FindCommitsByRepositoryID: %w", err)
	}
	defer rows.Close()

	var commits []models.SubmissionCommit
	for rows.Next() {
		var commit models.SubmissionCommit
		if err := rows.Scan(&commit.ID, &commit.SubmissionID, &commit.RepositoryID,
			&commit.FilePath, &commit.CommitSha, &commit.CommitURL, &commit.CreatedAt); err != nil {
			return nil, fmt.Errorf("repository.FindCommitsByRepositoryID: %w", err)
		}
		commits = append(commits, commit)
	}
	return commits, rows.Err()
}
