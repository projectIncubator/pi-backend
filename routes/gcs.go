package routes

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/db/cloud"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func (app *App) RegisterGoogleCloudRoutes() {
	app.router.HandleFunc("/bucket/upload", app.AddObject).Methods("PUT")
	app.router.HandleFunc("/projects/{project_id}/delete/{filename}", app.DeleteObject).Methods("DELETE")

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
	// TODO validate that if the user is the admin of the project project.photos/project.logo
	type UserToken struct {
		ID        string `json:"id"`
	}

	var userTk UserToken
	reqBody, err := ioutil.ReadAll(r.Body) // Read the request body (an userID for now)
	if err != nil {
		log.Printf("App.DeleteObject - error reading request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &userTk)
	if err != nil {
		log.Printf("App.DeleteObject - error unmarshaling request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID :=userTk.ID
	projectID := mux.Vars(r)["project_id"]
	//Call a function to check if the user is an admin if the user is an admin
	isAdmin :=app.store.ProjectProvider.CheckAdmin(projectID, userID)
	if !isAdmin {
		log.Printf("App.DeleteObject - not an admin, do not have permission to delete %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fileName := mux.Vars(r)["filename"]
	err = cloud.GCSDelete(fileName)
	if err != nil {
		log.Printf("App.DeleteObject - internal error, can't delete %v", err)
		return
	}
}