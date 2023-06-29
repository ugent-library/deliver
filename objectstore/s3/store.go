package s3

import (
	"context"
	"io"
	"net/url"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ugent-library/deliver/objectstore"
)

func init() {
	objectstore.Register("s3", New)
}

// connection string format: http(s)://id:secret@endpoint/bucket?region=region
// see https://stackoverflow.com/questions/67575681/is-aws-go-sdk-v2-integrated-with-local-minio-server
func New(conn string) (objectstore.Store, error) {
	u, err := url.Parse(conn)
	if err != nil {
		return nil, err
	}

	endpoint := u.Scheme + "://" + u.Host
	bucket := u.Path[1:]
	region := u.Query().Get("region")
	id := u.User.Username()
	secret, _ := u.User.Password()

	config := aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentialsProvider(id, secret, ""),
	}
	if u.Host != "" {
		config.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, region string, opts ...any) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               endpoint,
				SigningRegion:     region,
				HostnameImmutable: true,
			}, nil
		})
	}

	return &s3storage{
		client: s3.NewFromConfig(config),
		bucket: bucket,
	}, nil
}

type s3storage struct {
	client *s3.Client
	bucket string
}

func (s *s3storage) Add(ctx context.Context, id string, b io.Reader) (string, error) {
	uploader := manager.NewUploader(s.client)
	res, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(id),
		Body:   b,
	})
	if err != nil {
		return "", err
	}
	md5, _ := strconv.Unquote(*res.ETag)
	return md5, nil
}

func (s *s3storage) Get(ctx context.Context, id string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(id),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

func (s *s3storage) Delete(ctx context.Context, id string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(id),
	})
	return err
}

func (s *s3storage) IterateID(ctx context.Context, fn func(string) error) error {
	pager := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	})

	for pager.HasMorePages() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, obj := range page.Contents {
			err := fn(*obj.Key)
			if err == objectstore.Stop {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}
