package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"pr-reviewer-2025/internal/models"
)

type Teams interface {
	CreateOrUpdateTeam(ctx context.Context, team models.Team) (models.Team, error)
	GetTeam(ctx context.Context, teamName string) (models.Team, error)
}

type Users interface {
	SetUserActive(ctx context.Context, userID string, isActive bool) (models.User, error)
	GetUserByID(ctx context.Context, userID string) (models.User, error)
}

type PullRequests interface {
	CreatePR(ctx context.Context, req models.CreatePRRequest) (models.PullRequest, error)
	MergePR(ctx context.Context, prID string) (models.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (models.PullRequest, string, error)
	GetPRsByReviewer(ctx context.Context, reviewerID string) ([]models.PullRequestShort, error)
}

type Repository struct {
	Teams
	Users
	PullRequests
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Teams:        NewTeamsPostgres(db),
		Users:        NewUsersPostgres(db),
		PullRequests: NewPullRequestsPostgres(db),
	}
}
