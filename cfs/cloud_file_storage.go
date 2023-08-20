package cfs

import (
	"net/http"
	"time"
)

type CloudFileStorage interface {
	PutObject(bucket string, key string, file string, contentType string) error
	GetObject(bucket string, key string, file string) error
	DeleteObject(bucket string, key string) error
	DeleteObjects(bucket string, keys []string) error
	SignPutObjectURL(bucket string, key string, dur time.Duration, contentType string) (*PresignedURL, error)
	SignGetObjectURL(bucket string, key string, dur time.Duration) (string, error)
}

type PresignedURL struct {
	URL          string      `json:"url"`
	Method       string      `json:"method"`
	SignedHeader http.Header `json:"signed_header"`
}
