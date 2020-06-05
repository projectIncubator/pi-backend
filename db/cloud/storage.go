package cloud

import (
	"cloud.google.com/go/storage"
	"log"
	"time"
	"os"
	"context"
)

func GCSDelete (name string) error {
	bucketName := os.Getenv("BUCKET_NAME")
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	client, _ := storage.NewClient(ctx)


	o := client.Bucket(bucketName).Object(name)
	if err := o.Delete(ctx); err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}
}