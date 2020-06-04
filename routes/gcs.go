package routes

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func (app *App) RegisterGoogleCloudRoutes() {
	app.router.HandleFunc("/bucket/upload", app.AddObject).Methods("PUT")
	app.router.HandleFunc("/bucket/delete/{objectName}", app.DeleteObject).Methods("DELETE")

}

func (app *App) AddObject(w http.ResponseWriter, r *http.Request) {
	log.Printf("testing entering")
	//projectID := os.Getenv("YOUR_PROJECT_ID")
	bucketName := os.Getenv("BUCKET_NAME")

	// Opening file
	f, err := os.Open("notes.txt")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer f.Close()

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// bucket := client.Bucket(bucketName)
	object := "johnsnotes.txt"
	wc := client.Bucket(bucketName).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := wc.Close(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}


func (app *App) DeleteObject(w http.ResponseWriter, r *http.Request) {
	bucketName := os.Getenv("BUCKET_NAME")
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	client, _ := storage.NewClient(ctx)
	object := mux.Vars(r)["objectName"]

	o := client.Bucket(bucketName).Object(object)
	if err := o.Delete(ctx); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	} else {
		return
	}
}