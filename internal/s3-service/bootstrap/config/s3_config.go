package config

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	AccessKey, SecretKey, Url string
}

func NewS3Client(
	s3Config *S3Config,
) (*s3.Client, error) {
	cfg := aws.Config{}
	cfg.Region = "eu-central-2"
	cfg.Credentials = credentials.NewStaticCredentialsProvider(
		s3Config.AccessKey,
		s3Config.SecretKey,
		"",
	) 
	cfg.BaseEndpoint = aws.String(s3Config.Url)
	
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return s3Client, nil
}
