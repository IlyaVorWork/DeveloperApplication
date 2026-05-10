package s3

import (
	"context"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileStorage interface {
	UploadFile(ctx context.Context, fileName string, reader io.Reader, size int64, contentType string) error
	GetFileLink(fileName string) (string, error)
}

type Minio struct {
	client        *minio.Client
	presignClient *minio.Client
	bucket        string
}

func NewMinio(endpoint, access, secret, bucket, publicEndpoint, region string) *Minio {
	newClient := func(ep string) *minio.Client {
		c, err := minio.New(ep, &minio.Options{
			Creds:  credentials.NewStaticV4(access, secret, ""),
			Region: region,
		})
		if err != nil {
			log.Fatalln(err)
		}
		return c
	}

	main := newClient(endpoint)
	presign := main
	if publicEndpoint != "" && publicEndpoint != endpoint {
		presign = newClient(publicEndpoint)
	}

	return &Minio{
		client:        main,
		presignClient: presign,
		bucket:        bucket,
	}
}

func (m *Minio) UploadFile(ctx context.Context, fileName string, reader io.Reader, size int64, contentType string) error {
	_, err := m.client.PutObject(ctx, m.bucket, fileName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (m *Minio) GetFileLink(fileName string) (string, error) {
	ctx := context.Background()

	u, err := m.presignClient.PresignedGetObject(ctx, m.bucket, fileName, 15*time.Minute, url.Values{})
	if err != nil {
		return "", err
	}

	return u.String(), nil
}
