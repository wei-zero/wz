package s3cfs

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestS3Service_PutObject(t *testing.T) {
	godotenv.Load()
	s3Service, err := New()
	assert.NoError(t, err)
	err = s3Service.PutObject("xiaoma-dev", "test/1_png", "pkg/img/water/testdata/1.png", "image/png")
	assert.NoError(t, err)
}

func TestSignPutURL(t *testing.T) {
	godotenv.Load()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	assert.NoError(t, err)

	s3client := s3.NewFromConfig(cfg)
	svc := s3.NewPresignClient(s3client)
	input := &s3.PutObjectInput{
		Bucket: aws.String("xiaoma1"),
		Key:    aws.String("test/1.png"),
		ACL:    types.ObjectCannedACLPublicRead,
	}

	preReq, err := svc.PresignPutObject(context.TODO(), input, s3.WithPresignExpires(time.Minute*30))
	assert.NoError(t, err)
	t.Log(preReq.URL)
	t.Log(preReq.SignedHeader)

	// put object
	f, err := os.Open("pkg/img/water/testdata/1.png")
	assert.NoError(t, err)
	st, err := f.Stat()
	assert.NoError(t, err)
	req, err := http.NewRequest(http.MethodPut, preReq.URL, nil)
	assert.NoError(t, err)
	for k, vs := range preReq.SignedHeader {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	req.Header.Add("Content-Type", "image/png")
	req.ContentLength = st.Size()
	req.Body = f

	t.Log(req.Header)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	t.Log(resp.StatusCode)
	b, _ := io.ReadAll(resp.Body)
	t.Log(string(b))

	defer f.Close()
}

func TestPresignGetURL(t *testing.T) {
	godotenv.Load()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	assert.NoError(t, err)

	s3client := s3.NewFromConfig(cfg)
	svc := s3.NewPresignClient(s3client)
	req, err := svc.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("xiaoma1"),
		Key:    aws.String("test/1.png"),
	}, s3.WithPresignExpires(time.Minute*30))
	assert.NoError(t, err)
	t.Log(req.URL)
}

func BenchmarkPresign(b *testing.B) {
	godotenv.Load()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	assert.NoError(b, err)

	s3client := s3.NewFromConfig(cfg)
	svc := s3.NewPresignClient(s3client)

	for i := 0; i < b.N; i++ {
		_, err := svc.PresignGetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String("xiaoma1"),
			Key:    aws.String("test/1.png"),
		}, s3.WithPresignExpires(time.Minute*30))
		assert.NoError(b, err)
	}
}
