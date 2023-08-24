package mockcfs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/wei-zero/wz/cfs"
)

type mockCloudFileStorage struct {
	root string
}

// DeleteObject implements cfs.CloudFileStorage.
func (s *mockCloudFileStorage) DeleteObject(bucket string, key string) error {
	return os.Remove(fmt.Sprintf("%s/%s/%s", s.root, bucket, key))
}

// DeleteObjects implements cfs.CloudFileStorage.
func (s *mockCloudFileStorage) DeleteObjects(bucket string, keys []string) error {
	for _, key := range keys {
		if err := os.Remove(fmt.Sprintf("%s/%s/%s", s.root, bucket, key)); err != nil {
			return err
		}
	}
	return nil
}

// GetObject implements cfs.CloudFileStorage.
func (s *mockCloudFileStorage) GetObject(bucket string, key string, file string) error {
	dstFilePath := filepath.Join(s.root, bucket, key)
	source, err := os.Open(dstFilePath)
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
func (s *mockCloudFileStorage) PutObject(bucket string, key string, file string, contentType string) error {
	source, err := os.Open(file)
	if err != nil {
		return err
	}

	dstFilePath := filepath.Join(s.root, bucket, key)
	dstDir := filepath.Dir(dstFilePath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}
	dest, err := os.Create(dstFilePath)
	if err != nil {
		return err
	}

	// copy file
	_, err = io.Copy(dest, source)
	return err
}

// SignGetObjectURL implements cfs.CloudFileStorage.
func (s *mockCloudFileStorage) SignGetObjectURL(bucket string, key string, dur time.Duration) (string, error) {
	dstFilePath := filepath.Join(s.root, bucket, key)
	return dstFilePath, nil
}

// SignPutObjectURL implements cfs.CloudFileStorage.
func (s *mockCloudFileStorage) SignPutObjectURL(bucket string, key string, dur time.Duration, contentType string) (*cfs.PresignedURL, error) {
	panic("unimplemented")
}

func New() (cfs.CloudFileStorage, error) {
	root := os.Expand("$HOME/.mockcfs", os.Getenv)
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}
	return &mockCloudFileStorage{root: root}, nil
}
