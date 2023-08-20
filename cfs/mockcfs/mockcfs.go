package mockcfs

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/wei-zero/wz/cfs"
)

type mockCloudFileStorage struct {
}

// DeleteObject implements cfs.CloudFileStorage.
func (*mockCloudFileStorage) DeleteObject(bucket string, key string) error {
	return os.Remove(fmt.Sprintf("%s/%s", bucket, key))
}

// DeleteObjects implements cfs.CloudFileStorage.
func (*mockCloudFileStorage) DeleteObjects(bucket string, keys []string) error {
	for _, key := range keys {
		if err := os.Remove(fmt.Sprintf("%s/%s", bucket, key)); err != nil {
			return err
		}
	}
	return nil
}

// GetObject implements cfs.CloudFileStorage.
func (*mockCloudFileStorage) GetObject(bucket string, key string, file string) error {
	source, err := os.Open(fmt.Sprintf("%s/%s", bucket, key))
	if err != nil {
		return err
	}

	dest, err := os.Create(file)
	if err != nil {
		return err
	}

	// copy file
	_, err = io.Copy(dest, source)
	return err
}

// PutObject implements cfs.CloudFileStorage.
func (*mockCloudFileStorage) PutObject(bucket string, key string, file string, contentType string) error {
	source, err := os.Open(file)
	if err != nil {
		return err
	}

	dest, err := os.Create(fmt.Sprintf("%s/%s", bucket, key))
	if err != nil {
		return err
	}

	// copy file
	_, err = io.Copy(dest, source)
	return err
}

// SignGetObjectURL implements cfs.CloudFileStorage.
func (*mockCloudFileStorage) SignGetObjectURL(bucket string, key string, dur time.Duration) (string, error) {
	panic("unimplemented")
}

// SignPutObjectURL implements cfs.CloudFileStorage.
func (*mockCloudFileStorage) SignPutObjectURL(bucket string, key string, dur time.Duration, contentType string) (*cfs.PresignedURL, error) {
	panic("unimplemented")
}

func New() (cfs.CloudFileStorage, error) {
	return &mockCloudFileStorage{}, nil
}