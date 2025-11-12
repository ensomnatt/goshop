package usecase

import (
	"context"
	"errors"
	"time"
	"users_service/internal/config"
	"users_service/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UsersUseCase struct {
	repo domain.UserRepository
	cfg  *config.Config
}

type Claims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

func NewUsersUseCase(repo domain.UserRepository, cfg *config.Config) *UsersUseCase {
	return &UsersUseCase{
		repo: repo,
		cfg:  cfg,
	}
}

func (uc *UsersUseCase) Register(name, email, password string, ctx context.Context) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	err = uc.repo.Create(user, ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UsersUseCase) Login(email, password string, ctx context.Context) (string, time.Time, error) {
	expirationTime := time.Now().Add(time.Hour * 24)

	user, err := uc.repo.GetByEmail(email, ctx)
	if err != nil {
		return "", expirationTime, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", expirationTime, errors.New("Invalid credentials")
	}

	claims := &Claims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(uc.cfg.JWTSecret))
	if err != nil {
		return "", expirationTime, err
	}

	return tokenString, expirationTime, nil
}
