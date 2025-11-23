package service

import (
	"context"
	"pr-reviewer-2025/internal/models"
	"pr-reviewer-2025/internal/repository"
)

type TeamsService struct {
	repo repository.Teams
}

func NewTeamsService(repo repository.Teams) *TeamsService {
	return &TeamsService{repo: repo}
}

func (s *TeamsService) CreateOrUpdateTeam(ctx context.Context, team models.Team) (models.Team, error) {
	return s.repo.CreateOrUpdateTeam(ctx, team)
}

func (s *TeamsService) GetTeam(ctx context.Context, teamName string) (models.Team, error) {
	return s.repo.GetTeam(ctx, teamName)
}
