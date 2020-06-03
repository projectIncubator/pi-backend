package cloud

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"log"
	"mime/multipart"
	"time"
	"os"
)

func GCSUploader(name string, imageFile multipart.File) (error) {
	bucketName := os.Getenv("BUCKET_NAME")
	// Google Cloud Storage function
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return err
	}

	wc := client.Bucket(bucketName).Object(name).NewWriter(ctx)
	if _, err = io.Copy(wc, imageFile); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}