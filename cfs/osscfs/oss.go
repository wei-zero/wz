package osscfs

import (
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/wei-zero/wz/cfs"
)

type ossService struct {
	client *oss.Client
}

func New(endpoint string, accessKeyID string, accessKeySecret string) (cfs.CloudFileStorage, error) {

	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}
	return &ossService{
		client: client,
	}, nil
}

func (o *ossService) PutObject(bucket string, key string, file string, contentType string) error {
	b, err := o.client.Bucket(bucket)
	if err != nil {
		return err
	}

	return b.PutObjectFromFile(key, file, oss.ContentType(contentType))
}

func (o *ossService) GetObject(bucket string, key string, file string) error {
	b, err := o.client.Bucket(bucket)
	if err != nil {
		return err
	}

	return b.GetObjectToFile(key, file)
}

func (o *ossService) DeleteObject(bucket string, key string) error {
	b, err := o.client.Bucket(bucket)
	if err != nil {
		return err
	}

	return b.DeleteObject(key)
}

func (o *ossService) DeleteObjects(bucket string, keys []string) error {
	b, err := o.client.Bucket(bucket)
	if err != nil {
		return err
	}

	_, err = b.DeleteObjects(keys)
	return err
}

func (o *ossService) SignPutObjectURL(bucket string, key string, dur time.Duration, contentType string) (*cfs.PresignedURL, error) {
	b, err := o.client.Bucket(bucket)
	if err != nil {
		return nil, err
	}

	options := []oss.Option{
		//oss.Meta("myprop", "mypropval"),
		oss.ContentType(contentType),
	}

	url, err := b.SignURL(key, oss.HTTPPut, int64(dur.Seconds()), options...)
	if err != nil {
		return nil, err
	}

	return &cfs.PresignedURL{
		URL:    url,
		Method: "PUT",
	}, nil
}

func (o *ossService) SignGetObjectURL(bucket string, key string, dur time.Duration) (string, error) {
	b, err := o.client.Bucket(bucket)
	if err != nil {
		return "", err
	}

	return b.SignURL(key, oss.HTTPGet, int64(dur.Seconds()))
}
