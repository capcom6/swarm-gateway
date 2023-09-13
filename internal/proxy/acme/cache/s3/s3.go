package s3

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"golang.org/x/crypto/acme/autocert"
)

const (
	OptionEndpoint        = "endpoint"
	OptionAccessKeyId     = "accessKeyId"
	OptionSecretAccessKey = "secretAccessKey"
	OptionRegion          = "region"
	OptionBucket          = "bucket"
	OptionPrefix          = "prefix"
)

type S3Cache struct {
	Client *s3.S3

	Bucket string
	Prefix string
}

func New(options map[string]string) (*S3Cache, error) {
	cfg := aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(options[OptionAccessKeyId], options[OptionSecretAccessKey], "")).
		WithRegion(options[OptionRegion]).
		WithEndpoint(options[OptionEndpoint])
	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	return &S3Cache{
		Client: s3.New(sess),
		Bucket: options[OptionBucket],
		Prefix: options[OptionPrefix],
	}, nil
}

// Get returns a certificate data for the specified key.
// If there's no such key, Get returns ErrCacheMiss.
func (s *S3Cache) Get(ctx context.Context, key string) ([]byte, error) {
	object, err := s.Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s.Prefix + key),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchKey {
				return nil, autocert.ErrCacheMiss
			}
		}

		return nil, err
	}
	defer object.Body.Close()

	data, err := io.ReadAll(object.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Put stores the data in the cache under the specified key.
// Underlying implementations may use any data storage format,
// as long as the reverse operation, Get, results in the original data.
func (s *S3Cache) Put(ctx context.Context, key string, data []byte) error {
	_, err := s.Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s.Prefix + key),
		Body:   bytes.NewReader(data),
	})

	return err
}

// Delete removes a certificate data from the cache under the specified key.
// If there's no such key in the cache, Delete returns nil.
func (s *S3Cache) Delete(ctx context.Context, key string) error {
	_, err := s.Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s.Prefix + key),
	})

	return err
}
