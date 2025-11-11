package repository

import (
	"context"
	"errors"
	"users_service/internal/config"
	"users_service/internal/domain"
	"users_service/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	DB *pgxpool.Pool
}

func NewPostgresPool(cfg *config.Config, log *logger.Logger) *pgxpool.Pool {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Failed to ping db: %v", err)
	}

	log.Info("Connected to db")
	return pool
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &UserRepo{
		DB: db,
	}
}

func (r *UserRepo) Create(user *domain.User, ctx context.Context) error {
	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id"
	return r.DB.QueryRow(ctx, query, user.Name, user.Email, user.Password).Scan(&user.ID)
}

func (r *UserRepo) GetByEmail(email string, ctx context.Context) (*domain.User, error) {
	user := &domain.User{}
	query := "SELECT id, name, email, password, created_at FROM users WHERE email = $1"
	err := r.DB.QueryRow(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return user, nil
}
