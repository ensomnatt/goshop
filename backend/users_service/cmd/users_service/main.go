package main

import (
	"context"
	"fmt"
	"users_service/internal/config"
	"users_service/internal/domain"
	"users_service/internal/repository"
	"users_service/pkg/logger"
)

func main() {
	log := logger.New("debug")
	cfg := config.Load()
	pool := repository.NewPostgresPool(cfg, log)
	repo := repository.NewUserRepository(pool)

	user := &domain.User{
		Name:  "shushmyr",
		Email: "ensomnatt@protonmail.com",
	}

	ctx := context.Background()

	err := repo.Create(user, ctx)
	if err != nil {
		log.Errorf("Failed to create user: %v", err)
	}

	userFromDB, err := repo.GetByEmail(user.Email, ctx)
	if err != nil {
		log.Errorf("Failed to get user: %v", err)
	}

	fmt.Println(userFromDB)
}
