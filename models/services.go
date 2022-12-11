package models

import (
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	DatabaseURL       string
	S3Endpoint        string
	S3AccessKeyID     string
	S3SecretAccessKey string
	S3Bucket          string
}

type Services struct {
	Repository RepositoryService
	File       FileService
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
