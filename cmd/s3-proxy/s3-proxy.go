package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/istt/api_gateway/internal/app"
	minio "github.com/minio/minio-go/v7"
)

var configFile string

// main function
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.StringVar(&configFile, "config", "configs/api-gateway.yaml", "API Gateway configuration file")
	flag.Parse()

	// 1 - set default settings for components.
	app.S3Config()
	// 2 - override defaults with configuration file and watch changes
	app.ConfigInit(configFile)
	app.ConfigWatch(configFile)

	// 3 - bring up components
	app.S3Init()
	// 4 - setup the web server
	srv := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				return c.Status(code).JSON(e)
			}
			return c.Status(code).JSON(fiber.Map{"error": code, "message": err.Error()})
		},
	})
	setupRoutes(srv)

	log.Fatal(srv.Listen(app.Config.String("http.listen")))
}

// setupRoutes for file uploading
func setupRoutes(app *fiber.App) {
	// + eptw regions endpoints
	app.Post("api/file-upload", FileUpload)
	app.Post("api/file-upload/:name", FileUpload)
	app.Post("api/file-upload/:bucket/:name", FileUpload)

	app.Delete("api/file-upload/:name", FileRemove)
	app.Delete("api/file-upload/:bucket/:name", FileRemove)

	app.Get("api/statics/:name", FileDownload)
	app.Get("api/statics/:bucket/:name", FileDownload)
	app.Get("api/file-download", FileDownload)

	app.Get("api/file-info/:bucket/:name", FileCheck)
	app.Get("api/file-info/:name", FileCheck)
	app.Get("api/file-check", FileCheck)
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
		_, err := app.S3Client.StatObject(context.Background(), bucketName, newFileName, minio.StatObjectOptions{})
		if err != nil {

			uploadInfo, err := app.S3Client.PutObject(context.Background(), bucketName, newFileName, file, fileStat.Size(), minio.PutObjectOptions{ContentType: fileMime})
			if err == nil {
				return c.JSON(uploadInfo)
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
