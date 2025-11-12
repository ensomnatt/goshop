package main

import (
	"users_service/internal/config"
	httpDelievery "users_service/internal/delivery/http"
	"users_service/internal/repository"
	"users_service/internal/usecase"
	"users_service/pkg/logger"
)

func main() {
	log := logger.New("debug")
	cfg := config.Load()
	pool := repository.NewPostgresPool(cfg, log)
	repo := repository.NewUserRepository(pool)
	uc := usecase.NewUsersUseCase(repo, cfg)
	handler := httpDelievery.NewUsersHandler(uc, cfg, log)

	log.Infof("Server is up and running, port: %s", cfg.Port)
	err := handler.Start()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
