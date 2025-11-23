package service

import (
	"context"
	"pr-reviewer-2025/internal/models"
	"pr-reviewer-2025/internal/repository"
)

type PullRequestsService struct {
	repo repository.PullRequests
}

func NewPullRequestsService(repo repository.PullRequests) *PullRequestsService {
	return &PullRequestsService{repo: repo}
}

func (s *PullRequestsService) CreatePR(ctx context.Context, req models.CreatePRRequest) (models.PullRequest, error) {
	return s.repo.CreatePR(ctx, req)
}

func (s *PullRequestsService) MergePR(ctx context.Context, prID string) (models.PullRequest, error) {
	return s.repo.MergePR(ctx, prID)
}

func (s *PullRequestsService) ReassignReviewer(ctx context.Context, prID, oldUserID string) (models.PullRequest, string, error) {
	return s.repo.ReassignReviewer(ctx, prID, oldUserID)
}
