package s3storage

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/wei-zero/wz/cfs"
)

type s3Service struct {
	client        *s3.Client
	presignClient *s3.PresignClient
}

func New() (cfs.CloudFileStorage, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	s3client := s3.NewFromConfig(cfg)
	return &s3Service{
		client:        s3client,
		presignClient: s3.NewPresignClient(s3client),
	}, nil
}

func (s *s3Service) PutObject(bucket string, key string, filePath string, contentType string) error {
	stat, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if contentType == "" {
		contentType = "binary/octet-stream"
	}

	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		ACL:           types.ObjectCannedACLPublicRead,
		Body:          file,
		ContentLength: stat.Size(),
		ContentType:   aws.String(contentType),
	})

	return err
}

func (s *s3Service) GetObject(bucket string, key string, saveFilePath string) error {
	resp, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	file, err := os.Create(saveFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func (s *s3Service) DeleteObject(bucket string, key string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return err
}

func (s *s3Service) DeleteObjects(bucket string, keys []string) error {
	var deleteObjects []types.ObjectIdentifier
	for _, key := range keys {
		deleteObjects = append(deleteObjects, types.ObjectIdentifier{Key: aws.String(key)})
	}

	_, err := s.client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &types.Delete{Objects: deleteObjects},
	})

	return err
}

func (s *s3Service) SignPutObjectURL(bucket string, key string, dur time.Duration, contentType string) (*cfs.PresignedURL, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}
	//if publicRead {
	//	input.ACL = types.ObjectCannedACLPublicRead
	//}

	preReq, err := s.presignClient.PresignPutObject(context.TODO(), input, s3.WithPresignExpires(dur))
	if err != nil {
		return nil, err
	}

	return &cfs.PresignedURL{
		URL:          preReq.URL,
		Method:       preReq.Method,
		SignedHeader: preReq.SignedHeader,
	}, nil
}

func (s *s3Service) SignGetObjectURL(bucket string, key string, dur time.Duration) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	req, err := s.presignClient.PresignGetObject(context.TODO(), input, s3.WithPresignExpires(dur))
	if err != nil {
		return "", err
	}

	return req.URL, nil
}
