package service

import (
	"context"
	"pr-reviewer-2025/internal/models"
	"pr-reviewer-2025/internal/repository"
)

type Teams interface {
	CreateOrUpdateTeam(ctx context.Context, team models.Team) (models.Team, error)
	GetTeam(ctx context.Context, teamName string) (models.Team, error)
}

type Users interface {
	SetUserActive(ctx context.Context, userID string, isActive bool) (models.User, error)
	GetPRsForReviewer(ctx context.Context, reviewerID string) ([]models.PullRequestShort, error)
}

type PullRequests interface {
	CreatePR(ctx context.Context, req models.CreatePRRequest) (models.PullRequest, error)
	MergePR(ctx context.Context, prID string) (models.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (models.PullRequest, string, error)
}

type Service struct {
	Teams
	Users
	PullRequests
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Teams:        NewTeamsService(repos.Teams),
		Users:        NewUsersService(repos.Users, repos.PullRequests),
		PullRequests: NewPullRequestsService(repos.PullRequests),
	}
}
