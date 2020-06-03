package routes

import (
	"github.com/gorilla/mux"
	"go-api/db/cloud"
	"go-api/utils"
	"log"
	"net/http"
	"github.com/google/uuid"
)

func (app *App) RegisterGoogleCloudRoutes() {
	app.router.HandleFunc("/project/{project_id}/upload", app.AddObject).Methods("PUT").Queries("destination", "{destination}")
}

const MaxFileSize = 6 << 20 // 6 MB

func (app *App) AddObject(w http.ResponseWriter, r *http.Request) {
	// get url info
	projectID := mux.Vars(r)["project_id"]
	destination := r.FormValue("destination")

	// TODO: Check if projectID empty or valid
	// TODO: Check if action done by admins --> someone who is allowed

	// Read multipart form (image)
	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize+512)
	parseErr := r.ParseMultipartForm(MaxFileSize)
	if parseErr != nil {
		log.Println("App.AddObject - failed to parse multipart form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		log.Println("App.AddObject - expecting multipart form file")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	imageFile, header, err := r.FormFile("image")

	if err != nil {
		log.Println("App.AddObject - image is absent: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ext, err := utils.CheckMime(imageFile)

	name := "projects/"+projectID+"/"+uuid.New().String()+ext

	// TODO: Add filename to Media struct
	filename := header.Filename

	// TODO: If destination == coverphoto - set coverphoto
	// TODO: If destination == logo - set logo
	// TODO: If destination == media, append to media

	err = cloud.GCSUploader(name, imageFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		if destination == "coverphoto" {
			// TODO: Set Project.Coverphoto as this new url
		}
		if destination == "logo" {
			// TODO: Set Project.Logo as this new url
		}
		// TODO: Append to Project.media regardless
		return
	}
}