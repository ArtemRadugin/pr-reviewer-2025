package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"pr-reviewer-2025/internal/models"
)

type TeamsPostgres struct {
	db *sqlx.DB
}

func NewTeamsPostgres(db *sqlx.DB) *TeamsPostgres {
	return &TeamsPostgres{db: db}
}

var ErrTeamExists = errors.New("team already exists")

func (r *TeamsPostgres) CreateOrUpdateTeam(ctx context.Context, team models.Team) (models.Team, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return models.Team{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(ctx, `INSERT INTO `+teamsTable+` (team_name) VALUES ($1)`, team.TeamName)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return models.Team{}, ErrTeamExists
		}
		return models.Team{}, err
	}

	for _, m := range team.Members {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO `+usersTable+` (user_id, username, team_name, is_active) 
             VALUES ($1,$2,$3,$4)
             ON CONFLICT (user_id) DO UPDATE SET username = EXCLUDED.username, team_name = EXCLUDED.team_name, is_active = EXCLUDED.is_active`,
			m.UserID, m.Username, team.TeamName, m.IsActive,
		)
		if err != nil {
			return models.Team{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return models.Team{}, err
	}

	return team, nil
}

func (r *TeamsPostgres) GetTeam(ctx context.Context, teamName string) (models.Team, error) {
	var members []models.TeamMember
	query := `SELECT user_id, username, is_active FROM ` + usersTable + ` WHERE team_name = $1`
	if err := r.db.SelectContext(ctx, &members, query, teamName); err != nil {
		if err == sql.ErrNoRows || len(members) == 0 {
			return models.Team{}, fmt.Errorf("team not found")
		}
		return models.Team{}, err
	}
	if len(members) == 0 {
		return models.Team{}, fmt.Errorf("team not found")
	}
	return models.Team{TeamName: teamName, Members: members}, nil
}
