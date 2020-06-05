package routes

import (

	"github.com/gorilla/mux"
	"go-api/db/cloud"
	"go-api/utils"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/google/uuid"
)

func (app *App) RegisterGoogleCloudRoutes() {

	app.router.HandleFunc("/project/{project_id}/upload", app.AddObject).Methods("PUT").Queries("destination", "{destination}")

	app.router.HandleFunc("/projects/{project_id}/delete/{filename}", app.DeleteObject).Methods("DELETE")

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
	log.Println(filename)

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