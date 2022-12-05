package models

import (
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	DatabaseURL string
}

type Services struct {
	Repository RepositoryService
}

func NewServices(c Config) (*Services, error) {
	repository, err := NewRepositoryService(c)
	if err != nil {
		return nil, err
	}

	return &Services{
		Repository: repository,
	}, nil
}
