package models

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileService interface {
}

func NewFileService(c Config) (FileService, error) {
	client, err := minio.New(c.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.S3AccessKeyID, c.S3SecretAccessKey, ""),
		Secure: true, // TODO make configurable
	})
	if err != nil {
		return nil, err
	}
	return &fileService{
		client: client,
	}, nil
}

type fileService struct {
	client *minio.Client
}
