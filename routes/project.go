package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/model"
	"io/ioutil"
	"log"
	"net/http"
)

func (app *App) RegisterProjectRoutes() {
	app.router.HandleFunc("/projects", app.CreateProject).Methods("POST")
	app.router.HandleFunc("/projects/{id}", app.GetProject).Methods("GET")
	app.router.HandleFunc("/projects", app.UpdateProject).Methods("PATCH")
	app.router.HandleFunc("/projects/{id}", app.DeleteProject).Methods("DELETE") // TODO: We will not be deleting data. We will only put an account in a deactivated state

	app.router.HandleFunc("/projects/{id}/themes/{theme_name}", app.AddTheme).Methods("POST")
	app.router.HandleFunc("/projects/{id}/themes/{theme_name}", app.RemoveTheme).Methods("DELETE")

	app.router.HandleFunc("/projects/{proj_id}/members/{user_id}", app.DeleteMember).Methods("DELETE")
	app.router.HandleFunc("/projects/{proj_id}/members/{user_id}", app.ToggleAdmin).Methods("PATCH")

}

//When get project it should also return its members

func (app *App) CreateProject(w http.ResponseWriter, r *http.Request) {
	var newProject model.Project
	reqBody, err := ioutil.ReadAll(r.Body) // Read the request body
	if err != nil {
		log.Printf("App.CreateProject - error reading request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &newProject) // Fill newProject with the values coming from frontend
	if err != nil {
		log.Printf("App.CreateProject - error unmarshaling request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := app.store.ProjectProvider.CreateProject(&newProject)
	if err != nil {
		log.Printf("App.CreateProject - error creating project %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newProject.ID = id
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newProject)
	return
}

func (app *App) GetProject(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["id"]

	if projectID == "" {
		log.Printf("App.GetOneUser - empty project id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	project, err := app.store.ProjectProvider.GetProject(projectID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(project) // <- Sending the project as a json {id: ..., Title: ..., Stage ... , .. }
}

func (app *App) UpdateProject(w http.ResponseWriter, r *http.Request) {
	var updatedProject model.Project
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("App.UpdateProject - could not read r.Body with ioutil")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &updatedProject)
	if err != nil {
		log.Printf("App.UpdateProject - was unable to unmarshal changes")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// TODO: Validate that the updated project exists
	project, err := app.store.ProjectProvider.UpdateProject(&updatedProject)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(project) // <- Sending the project as a json {id: ..., Title: ..., Stage ... , .. }
}

func (app *App) DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["id"]

	if projectID == "" {
		log.Printf("App.RemoveProejct - empty project id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.ProjectProvider.RemoveProject(projectID)
	if err != nil {
		log.Printf("App.RemoveProject - error removing the project %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *App) AddTheme(w http.ResponseWriter, r *http.Request) {
	themeName := mux.Vars(r)["theme_name"]
	projectID := mux.Vars(r)["id"]

	if themeName == "" || projectID == "" {
		log.Printf("App.RemoveMember - empty project id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.ProjectProvider.AddTheme(themeName, projectID)

	if err != nil {
		log.Printf("App.CreateProject - error creating project %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	return
}

func (app *App) RemoveTheme(w http.ResponseWriter, r *http.Request) {
	themeName := mux.Vars(r)["theme_name"]
	projectID := mux.Vars(r)["id"]

	if themeName == "" || projectID == "" {
		log.Printf("App.RemoveProejct - empty project id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.ProjectProvider.RemoveTheme(themeName, projectID)
	if err != nil {
		log.Printf("App.RemoveProject - error removing the project %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)

	return
}

func (app *App) DeleteMember(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["proj_id"]
	userID := mux.Vars(r)["user_id"]

	if projectID == "" {
		log.Printf("App.RemoveMember - empty project id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if userID == "" {
		log.Printf("App.RemoveMember - empty user id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.ProjectProvider.RemoveMember(projectID, userID)
	if err != nil {
		log.Printf("App.RemoveProject - error removing the member %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *App) ToggleAdmin(w http.ResponseWriter, r *http.Request) {
	projectID := mux.Vars(r)["proj_id"]
	userID := mux.Vars(r)["user_id"]

	if projectID == "" {
		log.Printf("App.RemoveProejct - empty project id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if userID == "" {
		log.Printf("App.RemoveProejct - empty user id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.ProjectProvider.ChangeAdmin(projectID, userID)
	if err != nil {
		log.Printf("App.ChangeAdmin - error changing the admin %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
