package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"pr-reviewer-2025/internal/models"
	"time"
)

var (
	ErrPRNotFound  = errors.New("pr not found")
	ErrPRMerged    = errors.New("pr merged")
	ErrNotAssigned = errors.New("not assigned")
	ErrNoCandidate = errors.New("no candidate")
)

type prDB struct {
	PullRequestID   string         `db:"pull_request_id"`
	PullRequestName string         `db:"pull_request_name"`
	AuthorID        string         `db:"author_id"`
	Status          string         `db:"status"`
	CreatedAt       *time.Time     `db:"created_at"`
	MergedAt        *time.Time     `db:"merged_at"`
	AssignedRaw     pq.StringArray `db:"assigned_reviewers"`
}

type PullRequestsPostgres struct {
	db *sqlx.DB
}

func NewPullRequestsPostgres(db *sqlx.DB) *PullRequestsPostgres {
	return &PullRequestsPostgres{db: db}
}

func (r *PullRequestsPostgres) CreatePR(ctx context.Context, req models.CreatePRRequest) (models.PullRequest, error) {
	var tmp int
	err := r.db.GetContext(ctx, &tmp, `SELECT 1 FROM `+pullRequestsTable+` WHERE pull_request_id = $1`, req.PullRequestID)
	if err == nil {
		return models.PullRequest{}, fmt.Errorf("pr exists")
	}
	if err != nil && err != sql.ErrNoRows {
		return models.PullRequest{}, err
	}

	var author models.User
	if err := r.db.GetContext(ctx, &author, `SELECT user_id, username, team_name, is_active FROM `+usersTable+` WHERE user_id = $1`, req.AuthorID); err != nil {
		if err == sql.ErrNoRows {
			return models.PullRequest{}, fmt.Errorf("author not found")
		}
		return models.PullRequest{}, err
	}

	var candidates []string
	if err := r.db.SelectContext(ctx, &candidates, `SELECT user_id FROM `+usersTable+` WHERE team_name = $1 AND is_active = true AND user_id <> $2 LIMIT 2`, author.TeamName, author.UserID); err != nil {
		return models.PullRequest{}, err
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return models.PullRequest{}, err
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, `INSERT INTO `+pullRequestsTable+` (pull_request_id, pull_request_name, author_id, status, created_at) VALUES ($1,$2,$3,$4,$5)`,
		req.PullRequestID, req.PullRequestName, req.AuthorID, "OPEN", time.Now().UTC())
	if err != nil {
		return models.PullRequest{}, err
	}

	for _, c := range candidates {
		if _, err = tx.ExecContext(ctx, `INSERT INTO `+prReviewersTable+` (pull_request_id, reviewer_id) VALUES ($1, $2)`, req.PullRequestID, c); err != nil {
			return models.PullRequest{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		// commit failed — rollback will be attempted by deferred func
		return models.PullRequest{}, err
	}

	tx = nil

	var dbPR prDB
	query := `
SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status, pr.created_at, pr.merged_at,
       COALESCE(array_agg(r.reviewer_id ORDER BY r.reviewer_id) FILTER (WHERE r.reviewer_id IS NOT NULL), ARRAY[]::text[]) AS assigned_reviewers
FROM ` + pullRequestsTable + ` pr
LEFT JOIN ` + prReviewersTable + ` r ON r.pull_request_id = pr.pull_request_id
WHERE pr.pull_request_id = $1
GROUP BY pr.pull_request_id
`
	if err := r.db.GetContext(ctx, &dbPR, query, req.PullRequestID); err != nil {
		return models.PullRequest{}, err
	}

	pr := models.PullRequest{
		PullRequestID:     dbPR.PullRequestID,
		PullRequestName:   dbPR.PullRequestName,
		AuthorID:          dbPR.AuthorID,
		Status:            dbPR.Status,
		AssignedReviewers: []string(dbPR.AssignedRaw),
		CreatedAt:         dbPR.CreatedAt,
		MergedAt:          dbPR.MergedAt,
	}

	return pr, nil
}

func (r *PullRequestsPostgres) MergePR(ctx context.Context, prID string) (models.PullRequest, error) {
	res, err := r.db.ExecContext(ctx, `UPDATE `+pullRequestsTable+` SET status = 'MERGED', merged_at = $1 WHERE pull_request_id = $2`, time.Now().UTC(), prID)
	if err != nil {
		return models.PullRequest{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return models.PullRequest{}, fmt.Errorf("pr not found")
	}

	var dbPR prDB
	query := `
SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status, pr.created_at, pr.merged_at,
       COALESCE(array_agg(r.reviewer_id ORDER BY r.reviewer_id) FILTER (WHERE r.reviewer_id IS NOT NULL), ARRAY[]::text[]) AS assigned_reviewers
FROM ` + pullRequestsTable + ` pr
LEFT JOIN ` + prReviewersTable + ` r ON r.pull_request_id = pr.pull_request_id
WHERE pr.pull_request_id = $1
GROUP BY pr.pull_request_id
`
	if err := r.db.GetContext(ctx, &dbPR, query, prID); err != nil {
		return models.PullRequest{}, err
	}

	pr := models.PullRequest{
		PullRequestID:     dbPR.PullRequestID,
		PullRequestName:   dbPR.PullRequestName,
		AuthorID:          dbPR.AuthorID,
		Status:            dbPR.Status,
		AssignedReviewers: []string(dbPR.AssignedRaw),
		CreatedAt:         dbPR.CreatedAt,
		MergedAt:          dbPR.MergedAt,
	}

	return pr, nil
}

func (r *PullRequestsPostgres) ReassignReviewer(ctx context.Context, prID, oldUserID string) (models.PullRequest, string, error) {
	var status string
	if err := r.db.GetContext(ctx, &status, `SELECT status FROM `+pullRequestsTable+` WHERE pull_request_id = $1`, prID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PullRequest{}, "", ErrPRNotFound
		}
		return models.PullRequest{}, "", err
	}
	if status == "MERGED" {
		return models.PullRequest{}, "", ErrPRMerged
	}

	var isAssigned bool
	if err := r.db.GetContext(ctx, &isAssigned, `SELECT EXISTS(SELECT 1 FROM `+prReviewersTable+` WHERE pull_request_id = $1 AND reviewer_id = $2)`, prID, oldUserID); err != nil {
		return models.PullRequest{}, "", err
	}
	if !isAssigned {
		return models.PullRequest{}, "", ErrNotAssigned
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return models.PullRequest{}, "", err
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	var authorID, teamName string
	if err := tx.GetContext(ctx, &authorID, `SELECT author_id FROM `+pullRequestsTable+` WHERE pull_request_id = $1`, prID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PullRequest{}, "", ErrPRNotFound
		}
		return models.PullRequest{}, "", err
	}
	if err := tx.GetContext(ctx, &teamName, `SELECT team_name FROM `+usersTable+` WHERE user_id = $1`, authorID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PullRequest{}, "", ErrNoCandidate
		}
		return models.PullRequest{}, "", err
	}

	var candidate string
	candidateQuery := `
SELECT u.user_id
FROM ` + usersTable + ` u
WHERE u.team_name = $1
  AND u.is_active = true
  AND u.user_id <> $3
  AND u.user_id NOT IN (
    SELECT reviewer_id FROM ` + prReviewersTable + ` WHERE pull_request_id = $2
  )
ORDER BY u.user_id
LIMIT 1
FOR UPDATE SKIP LOCKED
`
	if err := tx.GetContext(ctx, &candidate, candidateQuery, teamName, prID, authorID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PullRequest{}, "", ErrNoCandidate
		}
		return models.PullRequest{}, "", err
	}
	if candidate == "" || candidate == oldUserID {
		return models.PullRequest{}, "", ErrNoCandidate
	}

	if _, err := tx.ExecContext(ctx, `INSERT INTO `+prReviewersTable+` (pull_request_id, reviewer_id) VALUES ($1,$2)`, prID, candidate); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return models.PullRequest{}, "", ErrNoCandidate
		}
		return models.PullRequest{}, "", err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM `+prReviewersTable+` WHERE pull_request_id = $1 AND reviewer_id = $2`, prID, oldUserID); err != nil {
		return models.PullRequest{}, "", err
	}

	if err := tx.Commit(); err != nil {
		return models.PullRequest{}, "", err
	}
	committed = true

	var pr models.PullRequest
	query := `
SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status, pr.created_at, pr.merged_at,
       COALESCE(array_agg(r.reviewer_id ORDER BY r.reviewer_id) FILTER (WHERE r.reviewer_id IS NOT NULL), ARRAY[]::text[]) AS assigned_reviewers
FROM ` + pullRequestsTable + ` pr
LEFT JOIN ` + prReviewersTable + ` r ON r.pull_request_id = pr.pull_request_id
WHERE pr.pull_request_id = $1
GROUP BY pr.pull_request_id
`
	if err := r.db.GetContext(ctx, &pr, query, prID); err != nil {
		return models.PullRequest{}, "", err
	}

	return pr, candidate, nil
}

func (r *PullRequestsPostgres) GetPRsByReviewer(ctx context.Context, reviewerID string) ([]models.PullRequestShort, error) {
	prs := []models.PullRequestShort{}

	query := `SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
              FROM ` + pullRequestsTable + ` pr
              JOIN ` + prReviewersTable + ` r ON r.pull_request_id = pr.pull_request_id
              WHERE r.reviewer_id = $1`

	if err := r.db.SelectContext(ctx, &prs, query, reviewerID); err != nil {
		if err == sql.ErrNoRows {
			return prs, nil
		}
		return nil, err
	}
	return prs, nil
}
