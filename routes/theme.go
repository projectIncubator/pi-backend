package routes

import (
	"encoding/json"
	"go-api/model"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *App) RegisterThemeRoutes() {
	app.router.HandleFunc("/themes", app.CreateTheme).Methods("POST")
	app.router.HandleFunc("/themes/{theme_name}", app.GetTheme).Methods("GET")
	app.router.HandleFunc("/themes", app.UpdateTheme).Methods("PATCH")

	//app.router.HandleFunc("/themes/{theme_name}", app.GetProjectsWithTheme).Methods("GET")

	app.router.HandleFunc("/themes/{theme_name}", app.DeleteTheme).Methods("DELETE")
}

func (app *App) CreateTheme(w http.ResponseWriter, r *http.Request) {
	newTheme := model.NewTheme()
	reqBody, err := ioutil.ReadAll(r.Body) // Read the request body
	if err != nil {
		log.Printf("App.CreateProject - error reading request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &newTheme) // Fill newTheme with the values coming from frontend
	if err != nil {
		log.Printf("App.CreateProject - error unmarshaling request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = app.store.ThemeProvider.CreateTheme(&newTheme) // changed to err = from id, err :=
	if err != nil {
		log.Printf("App.CreateTheme - error creating theme %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	//json.NewEncoder(w).Encode(newTheme)
}
func (app *App) GetTheme(w http.ResponseWriter, r *http.Request) {
	themeName := mux.Vars(r)["theme_name"]

	if themeName == "" {
		log.Printf("App.GetOneUser - empty project id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	theme, err := app.store.ThemeProvider.GetTheme(themeName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(theme) // <- Sending the theme as a json {id: ..., Title: ..., Stage ... , .. }
}
func (app *App) UpdateTheme(w http.ResponseWriter, r *http.Request) {
	updatedTheme := model.NewTheme()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("App.UpdateProject - could not read r.Body with ioutil")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &updatedTheme)
	if err != nil {
		log.Printf("App.UpdateTheme - was unable to unmarshal changes")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name, err := app.store.ThemeProvider.UpdateTheme(&updatedTheme)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(name) // <- Sending the project as a json {id: ..., Title: ..., Stage ... , .. }
}

// TODO ?
//func (app *App) GetProjectsWithTheme(w http.ResponseWriter, r *http.Request) {
//	themeName := mux.Vars(r)["id"]
//
//	if themeName == "" {
//		log.Printf("App.GetProjectsWithTheme - empty theme name")
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//	err := app.store.ThemeProvider.GetProjectsWithTheme(themeName)
//	if err != nil {
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//	w.WriteHeader(http.StatusOK)
//	//json.NewEncoder(w).Encode(projects) // <- Sending the project as a json {id: ..., Title: ..., Stage ... , .. }
//}
func (app *App) DeleteTheme(w http.ResponseWriter, r *http.Request) {
	themeName := mux.Vars(r)["theme_name"]

	if themeName == "" {
		log.Printf("App.DeleteTheme - empty theme name")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.ThemeProvider.DeleteTheme(themeName)
	if err != nil {
		log.Printf("App.DeleteTheme - error removing the theme %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
