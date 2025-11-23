package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"pr-reviewer-2025/internal/models"
)

type UsersPostgres struct {
	db *sqlx.DB
}

func NewUsersPostgres(db *sqlx.DB) *UsersPostgres {
	return &UsersPostgres{db: db}
}

func (r *UsersPostgres) SetUserActive(ctx context.Context, userID string, isActive bool) (models.User, error) {
	res, err := r.db.ExecContext(ctx, `UPDATE `+usersTable+` SET is_active = $1 WHERE user_id = $2`, isActive, userID)
	if err != nil {
		return models.User{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return models.User{}, fmt.Errorf("user not found")
	}

	var user models.User
	if err := r.db.GetContext(ctx, &user, `SELECT user_id, username, team_name, is_active FROM `+usersTable+` WHERE user_id = $1`, userID); err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found")
		}
		return models.User{}, err
	}

	return user, nil
}

func (r *UsersPostgres) GetUserByID(ctx context.Context, userID string) (models.User, error) {
	var user models.User
	if err := r.db.GetContext(ctx, &user, `SELECT user_id, username, team_name, is_active FROM `+usersTable+` WHERE user_id = $1`, userID); err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}
