package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type UserRepository interface {
	NewTx(ctx context.Context) (pgx.Tx, error)
	Save(ctx context.Context, tx pgx.Tx, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	Db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) UserRepository {
	return &repository{Db: db}
}

func (r *repository) NewTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	return tx, nil
}

func (r *repository) Save(ctx context.Context, tx pgx.Tx, user *User) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO users (user_id, email, password, last_login, created)
		VALUES ($1, $2, $3, $4, $5)
	`, &user.UserID, &user.Email, &user.Password, &user.LastLogin, &user.Created)

	if err != nil {
		pgErr := err.(*pgconn.PgError)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUserAlreadyExists
		}

		log.Error().Err(err).Send()
		return err
	}

	return nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}

	err := r.Db.QueryRow(ctx, `
		SELECT * FROM users WHERE email = $1
	`, email).Scan(&user.UserID, &user.Email, &user.Password, &user.LastLogin, &user.Created)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		log.Error().Err(err).Send()
		return nil, err
	}

	return user, nil
}
