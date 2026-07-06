package repository

import (
	"context"
	"fmt"

	"github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, user domain.User) error {
	query := `INSERT INTO users	(id, name, email,password_hash, role, email_verified, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	_, err := r.db.Exec(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.Role, user.EmailVerified, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("Create user: %w", err)
	}

	return nil

}

func (r *Repository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `SELECT
    id,
    name,
    email,
    password_hash,
    role,
    email_verified,
    created_at,
    updated_at
	FROM users
	WHERE email = $1`

	row := r.db.QueryRow(ctx, query, email)

	var user domain.User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return domain.User{}, fmt.Errorf("find user by email: %w", err)
	}

	return user, nil

}

func (r *Repository) CreateRefreshToken(ctx context.Context, token domain.RefreshToken) error {
	query := `
	INSERT INTO refresh_tokens (
	id,
	user_id,
	token_hash,
	expires_at,
	created_at,
	revoked_at
	)
	VALUES ($1, $2, $3, $4, $5, $6);
	`
	_, err := r.db.Exec(ctx, query, token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.CreatedAt, token.RevokedAt)

	if err != nil {
		return err
	}

	return nil
}
