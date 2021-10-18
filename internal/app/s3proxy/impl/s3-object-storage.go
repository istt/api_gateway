package impl

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"

	"github.com/istt/api_gateway/internal/app/s3proxy/interfaces"
	"github.com/minio/minio-go/v7"
)

// S3ObjectStorage implements services.ObjectStorage
type S3ObjectStorage struct {
	minioClient *minio.Client
	bucket      string
}

func NewS3ObjectStorageService(minioClient *minio.Client, bucket string) interfaces.ObjectStorage {
	return &S3ObjectStorage{minioClient: minioClient, bucket: bucket}
}

// Put update the object with given key
func (s *S3ObjectStorage) Put(key string, content []byte) error {
	reader := bytes.NewReader(content)
	f, err := s.minioClient.PutObject(context.TODO(), s.bucket, key, reader, reader.Size(), minio.PutObjectOptions{})
	log.Printf("Successfully upload file: %s %s", f.Bucket, f.Key)
	return err
}

// Get return the value of the key
func (s *S3ObjectStorage) Get(key string) ([]byte, error) {
	obj, err := s.minioClient.GetObject(context.TODO(), s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(obj)
}

// Del remove the value of the key
func (s *S3ObjectStorage) Del(key string) error {
	return s.minioClient.RemoveObject(context.TODO(), s.bucket, key, minio.RemoveObjectOptions{})
}
