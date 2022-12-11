package models

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileService interface {
	Add(context.Context, string, io.ReadSeekCloser) error
}

// see https://stackoverflow.com/questions/67575681/is-aws-go-sdk-v2-integrated-with-local-minio-server
func NewFileService(c Config) (FileService, error) {
	config := aws.Config{
		Region:      c.S3Region,
		Credentials: credentials.NewStaticCredentialsProvider(c.S3AccessKeyID, c.S3SecretAccessKey, ""),
	}
	if c.S3URL != "" {
		config.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, opts ...any) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               c.S3URL,
				SigningRegion:     c.S3Region,
				HostnameImmutable: true,
			}, nil
		})
	}

	client := s3.NewFromConfig(config)

	return &fileService{
		client: client,
		bucket: c.S3Bucket,
	}, nil
}

type fileService struct {
	client *s3.Client
	bucket string
}

func (f *fileService) Add(ctx context.Context, id string, b io.ReadSeekCloser) error {
	return nil
}
