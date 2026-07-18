package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Aneeshie/ecommerce/internal/common/database"
	"github.com/Aneeshie/ecommerce/internal/identity"
	"github.com/Aneeshie/ecommerce/internal/identity/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db database.QueryExecutor
}

func NewRepository(db database.QueryExecutor) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, user domain.User) error {
	query := `INSERT INTO users	(id, name, email,password_hash, role, email_verified, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	_, err := r.db.Exec(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.Role, user.EmailVerified, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return identity.ErrEmailAlreadyExists
			}
		}
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
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, identity.ErrUserNotFound
		}
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

func (r *Repository) FindRefreshTokenByHash(ctx context.Context, hash string) (domain.RefreshToken, error) {
	query := `
	SELECT
	id,
	user_id,
	token_hash,
	expires_at,
	created_at,
	revoked_at
	FROM refresh_tokens
	WHERE token_hash = $1;
	`
	row := r.db.QueryRow(ctx, query, hash)

	var t domain.RefreshToken

	err := row.Scan(
		&t.ID,
		&t.UserID,
		&t.TokenHash,
		&t.ExpiresAt,
		&t.CreatedAt,
		&t.RevokedAt,
	)

	if err != nil {
		return domain.RefreshToken{}, err
	}

	return t, nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	query := `
	SELECT
		id,
		name,
		email,
		password_hash,
		role,
		email_verified,
		created_at,
		updated_at
	FROM users
	WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

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
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, identity.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("find by user id: %w", err)
	}

	return user, nil
}
