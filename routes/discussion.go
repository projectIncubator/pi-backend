package routes

import (
	"encoding/json"
	"go-api/model"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *App) RegisterDiscussionRoutes() {
	app.router.HandleFunc("/projects/{proj_id}/discussions", app.CreateDiscussion).Methods("POST")
	app.router.HandleFunc("/projects/{proj_id}/discussions/{disc_num}", app.GetDiscussion).Methods("GET")
}

func (app *App) CreateDiscussion(w http.ResponseWriter, r *http.Request) {
	newDiscussion := model.NewDiscussionIn()
	reqBody, err := ioutil.ReadAll(r.Body) // Read the request body
	if err != nil {
		log.Printf("App.CreateProject - error reading request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &newDiscussion) // Fill newProject with the values coming from frontend
	if err != nil {
		log.Printf("App.CreateProject - error unmarshaling request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	projectID := mux.Vars(r)["proj_id"]

	id, err := app.store.DiscussionProvider.CreateDiscussion(projectID, &newDiscussion)
	if err != nil {
		log.Printf("App.CreateProject - error creating project %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var newDiscussionNum = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newDiscussionNum)
}

func (app *App) GetDiscussion(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["proj_id"]
	discNum := mux.Vars(r)["disc_num"]

	if projectID == "" {
		log.Printf("App.GetOneUser - empty project id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	discussion, err := app.store.DiscussionProvider.GetDiscussion(projectID, discNum)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(discussion) // <- Sending the project as a json {id: ..., Title: ..., Stage ... , .. }
}
