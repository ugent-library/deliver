package models

import (
	"context"
	"io"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileService interface {
	Add(context.Context, string, io.ReadSeekCloser) (string, error)
	Get(context.Context, string, io.Writer) error
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

func (f *fileService) Add(ctx context.Context, id string, b io.ReadSeekCloser) (string, error) {
	uploader := manager.NewUploader(f.client)
	res, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    aws.String(id),
		Body:   b,
	})
	if err != nil {
		return "", err
	}
	md5, _ := strconv.Unquote(*res.ETag)
	return md5, nil
}

func (f *fileService) Get(ctx context.Context, id string, b io.Writer) error {
	downloader := manager.NewDownloader(f.client)
	downloader.Concurrency = 1
	_, err := downloader.Download(context.TODO(), fakeWriterAt{b}, &s3.GetObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    aws.String(id),
	})
	if err != nil {
		return err
	}
	return nil
}

// implement io.WriterAt for a plain io.Writer
// only works correctly if downloader.Concurrency = 1
type fakeWriterAt struct {
	w io.Writer
}

func (fw fakeWriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
	return fw.w.Write(p)
}
