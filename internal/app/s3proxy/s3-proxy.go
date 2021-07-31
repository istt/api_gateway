package s3proxy

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/istt/api_gateway/internal/app"
	"github.com/minio/minio-go/v7"
)

func SetupRoutes(app *fiber.App) {
	app.Get("api/file-browser", FileListing)
	app.Get("api/file-browser/:name", FileCheck)
	app.Post("api/file-browser", FileUpload)
	app.Delete("api/file-browser/:name", FileRemove)

	app.Post("api/file-upload", FileUpload)
	app.Post("api/file-upload/:name", FileUpload)
	app.Post("api/file-upload/:bucket/:name", FileUpload)

	app.Delete("api/file-upload/:name", FileRemove)
	app.Delete("api/file-upload/:bucket/:name", FileRemove)

	app.Get("api/statics/:name", FileDownload)
	app.Get("api/statics/:bucket/:name", FileDownload)
	app.Get("api/file-download/:name", FileDownload)
	app.Get("api/file-download/:bucket/:name", FileDownload)
	app.Get("api/file-download", FileDownload)

	app.Get("api/file-info/:bucket/:name", FileCheck)
	app.Get("api/file-info/:name", FileCheck)
	app.Get("api/file-check", FileCheck)
}

// FileListing list the files information from S3
func FileListing(c *fiber.Ctx) error {
	bucketName := c.Get("bucket", c.Query("bucket", app.Config.String("s3.bucket")))
	exists, err := app.S3Client.BucketExists(c.Context(), bucketName)
	if err != nil {
		return err
	}
	if !exists {
		return fiber.ErrNotFound
	}
	entities := make([]minio.ObjectInfo, 0)
	options := minio.ListObjectsOptions{
		Recursive: false,
	}
	if c.Query("size") != "" {
		size, err := strconv.Atoi(c.Query("size"))
		if err == nil {
			options.MaxKeys = size
		}
	}
	if c.Query("name") != "" {
		options.Prefix = c.Query("name")
	}
	for file := range app.S3Client.ListObjects(c.Context(), bucketName, options) {
		entities = append(entities, file)
	}
	return c.JSON(entities)
}

// FileUpload upload a file into minio bucket
func FileUpload(c *fiber.Ctx) error {
	bucketName := c.Get("bucket", app.Config.String("s3.bucket"))
	if exists, ok := app.S3Client.BucketExists(c.Context(), bucketName); ok == nil {
		if !exists {
			err := app.S3Client.MakeBucket(c.Context(), bucketName, minio.MakeBucketOptions{})
			if err != nil {
				return err
			}
		}
	} else {
		return fiber.ErrBadGateway
	}
	// Get first file from form field "document":
	uploadedFile, err := c.FormFile("file")
	if err != nil {
		return err
	}
	log.Printf("Got file header %+v", uploadedFile.Header)
	fileName := c.Params("name", c.Get("name", uploadedFile.Filename))
	filePath := os.TempDir() + "/" + fileName
	err = c.SaveFile(uploadedFile, filePath)
	if err != nil {
		return err
	}
	// + upload file to bucket
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	fileMime := uploadedFile.Header.Get("Content-Type")
	if fileMime == "" {
		mime, err := mimetype.DetectReader(file)
		if err != nil {
			fileMime = mime.String()
		}
	}

	// TODO: Check if object exists, slug the basepath, and append suffix
	ext := filepath.Ext(fileName)
	baseName := slug.Make(strings.TrimSuffix(filepath.Base(fileName), ext))
	suffix := 0
	newFileName := fmt.Sprintf("%s%s", baseName, ext)
	for {
		_, err := app.S3Client.StatObject(c.Context(), bucketName, newFileName, minio.StatObjectOptions{})
		if err != nil {

			uploadInfo, err := app.S3Client.PutObject(c.Context(), bucketName, newFileName, file, fileStat.Size(), minio.PutObjectOptions{ContentType: fileMime})
			if err == nil {
				if info, err := app.S3Client.StatObject(context.Background(), uploadInfo.Bucket, uploadInfo.Key, minio.StatObjectOptions{}); err == nil {
					return c.JSON(info)
				}
			}
		}
		suffix++
		newFileName = fmt.Sprintf("%s-%d%s", baseName, suffix, ext)
	}
}

// FileDownload download a file from minio
func FileDownload(c *fiber.Ctx) error {
	bucketName := c.Params("bucket", c.Get("bucket", app.Config.String("s3.bucket")))
	fileName := c.Params("name", c.Get("name"))
	if fileName == "" {
		return fiber.ErrBadRequest
	}
	obj, err := app.S3Client.GetObject(context.Background(), bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	fileStat, err := obj.Stat()
	if err != nil {
		return err
	}
	c.Set(fiber.HeaderContentType, fileStat.ContentType)
	return c.SendStream(obj, int(fileStat.Size))
}

// FileCheck check if one object exists
func FileCheck(c *fiber.Ctx) error {
	bucketName := c.Params("bucket", c.Get("bucket", app.Config.String("s3.bucket")))
	fileName := c.Params("name", c.Get("name"))
	if fileName == "" {
		return fiber.ErrBadRequest
	}
	objInfo, err := app.S3Client.StatObject(context.Background(), bucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	return c.JSON(objInfo)
}

// FileRemove remove one file from specified bucket
func FileRemove(c *fiber.Ctx) error {
	bucketName := c.Params("bucket", c.Get("bucket", app.Config.String("s3.bucket")))
	fileName := c.Params("name", c.Get("name"))
	if fileName == "" {
		return fiber.ErrBadRequest
	}
	err := app.S3Client.RemoveObject(c.Context(), bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
