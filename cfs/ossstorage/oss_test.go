package ossstorage

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestOssService_PutObject(t *testing.T) {
	godotenv.Load()

	o, err := New(os.Getenv("OSS_ENDPOINT"), os.Getenv("OSS_ACCESS_KEY_ID"), os.Getenv("OSS_ACCESS_SECRET"))
	assert.NoError(t, err)
	err = o.PutObject(os.Getenv("OSS_BUCKET"), "test/1_png", "pkg/img/water/testdata/1.png", "image/png")
	assert.NoError(t, err)
}

func TestOssService_GetObject(t *testing.T) {
	godotenv.Load()

	o, err := New(os.Getenv("ALIYUN_OSS_ENDPOINT"), os.Getenv("ALIYUN_ACCESS_KEY_ID"), os.Getenv("ALIYUN_ACCESS_KEY_SECRET"))
	assert.NoError(t, err)
	url, err := o.SignGetObjectURL(os.Getenv("BUCKET"), "test/test_1_png.png", 30*time.Minute)
	assert.NoError(t, err)
	t.Log(url)
}
