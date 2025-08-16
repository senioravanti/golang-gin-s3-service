package handlers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BucketHandler struct {
	S3Client *s3.Client
}

func (bh *BucketHandler) CreateBucket(
	ctx context.Context,
	bucketName string,
) (*s3.CreateBucketOutput, error) {
	output, err := bh.S3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

