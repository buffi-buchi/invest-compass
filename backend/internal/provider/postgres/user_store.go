package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

var (
	//go:embed queries/create_user.sql
	createUserQuery string

	//go:embed queries/get_user_by_id.sql
	getUserByIDQuery string

	//go:embed queries/get_user_by_email.sql
	getUserByEmailQuery string
)

type UserStore struct {
	db  *pgxpool.Pool
	id  func() (uuid.UUID, error)
	now func() time.Time
}

func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{
		db:  db,
		id:  func() (uuid.UUID, error) { return uuid.NewV7() },
		now: func() time.Time { return time.Now().UTC() },
	}
}

func (s *UserStore) Create(ctx context.Context, user model.User) (model.User, error) {
	id, err := s.id()
	if err != nil {
		return model.User{}, fmt.Errorf("create user ID: %w", err)
	}

	user.ID = id
	user.CreateTime = s.now()

	_, err = s.db.Exec(ctx, createUserQuery, user.ID, user.Email, user.Password, user.CreateTime)
	if err != nil {
		// TODO: Add a function to check postgres errors.
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return model.User{}, errors.Join(errors.New("user already exists"), model.ErrAlreadyExists)
			}
		}
		return model.User{}, fmt.Errorf("insert user: %w", err)
	}

	return user, nil
}

func (s *UserStore) GetByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	row := s.db.QueryRow(ctx, getUserByIDQuery, id)

	var user model.User

	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreateTime)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, model.ErrNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("select user by ID: %w", err)
	}

	return user, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (model.User, error) {
	row := s.db.QueryRow(ctx, getUserByEmailQuery, email)

	var user model.User

	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreateTime)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, model.ErrNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("select user by email: %w", err)
	}

	return user, nil
}
