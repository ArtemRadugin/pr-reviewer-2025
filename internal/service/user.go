package service

import (
	"context"
	"pr-reviewer-2025/internal/models"
	"pr-reviewer-2025/internal/repository"
)

type UsersService struct {
	usersRepo repository.Users
	prRepo    repository.PullRequests
}

func NewUsersService(usersRepo repository.Users, prRepo repository.PullRequests) *UsersService {
	return &UsersService{
		usersRepo: usersRepo,
		prRepo:    prRepo,
	}
}

func (s *UsersService) SetUserActive(ctx context.Context, userID string, isActive bool) (models.User, error) {
	return s.usersRepo.SetUserActive(ctx, userID, isActive)
}

func (s *UsersService) GetPRsForReviewer(ctx context.Context, reviewerID string) ([]models.PullRequestShort, error) {
	return s.prRepo.GetPRsByReviewer(ctx, reviewerID)
}
