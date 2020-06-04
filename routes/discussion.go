package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/model"
	"io/ioutil"
	"log"
	"net/http"
)


func (app *App) RegisterDiscussionRoutes() {
	app.router.HandleFunc("/projects/{proj_id}/discussions", app.CreateDiscussion).Methods("POST")
	app.router.HandleFunc("/projects/{proj_id}/discussions/{disc_num}", app.GetDiscussion).Methods("GET")
	//app.router.HandleFunc("/project/{proj_id}/discussion", app.GetDiscussions).Methods("GET")
}

func (app *App) CreateDiscussion(w http.ResponseWriter, r *http.Request) {
	var newDiscussion model.DiscussionIn
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
	_ = json.NewEncoder(w).Encode(newDiscussionNum)
	return
}

//func (app *App) GetDiscussions(w http.ResponseWriter, r *http.Request) {
//	projectID := mux.Var(r)["proj_id"]
//
//	if projectID == "" {
//		log.Printf("App.GetOneUser - empty project id")
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	project, err := app.store.DiscussionProvider.GetDiscussions(projectID)
//	if err != nil {
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	_ = json.NewEncoder(w).Encode(project) // <- Sending the project as a json {id: ..., Title: ..., Stage ... , .. }
//}


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
	_ = json.NewEncoder(w).Encode(discussion) // <- Sending the project as a json {id: ..., Title: ..., Stage ... , .. }
}