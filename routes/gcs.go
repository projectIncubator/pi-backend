package routes

import (
	"encoding/json"
	"go-api/db/cloud"
	"go-api/utils"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (app *App) RegisterGoogleCloudRoutes() {

	app.router.HandleFunc("/project/{project_id}/upload/{destination}", app.AddObject).Methods("PUT")

	app.router.HandleFunc("/project/{project_id}/delete", app.DeleteObject).Methods("DELETE")

}

const MaxFileSize = 6 << 20 // 6 MB

func (app *App) AddObject(w http.ResponseWriter, r *http.Request) {
	// Init
	// todo: change this to os.getenv for baseurl
	baseUrl := "https://storage.cloud.com/graceful-castle-276900.appspot.com/"
	// get url info
	projectID := mux.Vars(r)["project_id"]
	destination := mux.Vars(r)["destination"]

	project, _ := app.store.ProjectProvider.GetProject(projectID)

	if project == nil { // Check that the project exists (not storing for some uuid that isnt a real project
		log.Println("App.AddObject - No existing project with this projectID")
		w.WriteHeader(http.StatusNotFound)
		return
	}

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
	ext, _ := utils.CheckMime(imageFile)

	name := "projects/" + projectID + "/" + uuid.New().String() + ext

	filename := header.Filename
	log.Println(filename)
	err = cloud.GCSUploader(name, imageFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		if destination == "cover_photo" {
			app.store.ProjectProvider.UpdateCoverPhoto(projectID, baseUrl+name)
		}
		if destination == "logo" {
			app.store.ProjectProvider.UpdateLogo(projectID, baseUrl+name)
		}
		app.store.ProjectProvider.CreateProjectMedia(projectID, baseUrl+name)
		return
	}
}

func (app *App) DeleteObject(w http.ResponseWriter, r *http.Request) {
	type UserToken struct {
		ID       string `json:"id"`
		FileName string `json:"name"`
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
	userID := userTk.ID
	fileName := userTk.FileName
	projectID := mux.Vars(r)["project_id"]
	//Call a function to check if the user is an admin if the user is an admin
	isAdmin := app.store.ProjectProvider.CheckAdmin(projectID, userID)
	if !isAdmin {
		log.Printf("App.DeleteObject - not an admin, do not have permission to delete %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = cloud.GCSDelete(fileName)
	if err != nil {
		log.Printf("App.DeleteObject - internal error, can't delete %v", err)
		return
	}
}
