package app

import (
	"context"
	"log"

	"github.com/istt/api_gateway/internal/app/s3proxy/impl"
	"github.com/istt/api_gateway/internal/app/s3proxy/interfaces"
	"github.com/knadh/koanf/providers/confmap"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	// S3Client hold the connection to database
	S3Client *minio.Client

	// S3Storage hold the storage for reuse
	S3Storage interfaces.ObjectStorage
)

// S3Config configure application runtime
func S3Config() {
	// koanf defautl values
	Config.Load(confmap.Provider(map[string]interface{}{
		"s3.endpoint":        "127.0.0.1:8333",
		"s3.accessKeyID":     "",
		"s3.secretAccessKey": "",
		"s3.useSSL":          false,
		"s3.bucket":          "",
	}, "."), nil)
}

// S3Init initiate database
func S3Init() {
	minioClient, err := minio.New(Config.String("s3.endpoint"), &minio.Options{
		Creds:  credentials.NewStaticV4(Config.String("s3.accessKeyID"), Config.String("s3.secretAccessKey"), ""),
		Secure: Config.Bool("s3.useSSL"),
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to S3 URL %s", minioClient.EndpointURL())
	bucketName := Config.String("s3.bucket")
	if exists, err := minioClient.BucketExists(context.TODO(), bucketName); err == nil {
		if !exists {
			err := minioClient.MakeBucket(context.TODO(), bucketName, minio.MakeBucketOptions{})
			if err != nil {
				log.Fatalf("Error creating bucket [%s] : %s", bucketName, err)
			}
		} else {
			log.Printf("Bucket name : %s", bucketName)
		}
	} else {
		log.Fatalf("Error checking bucket: %s", err)
	}
	S3Client = minioClient
	S3Storage = impl.NewS3ObjectStorageService(minioClient, bucketName)
}
