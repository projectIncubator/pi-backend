package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/model"
	"io/ioutil"
	"log"
	"net/http"
)

func (app *App) RegisterUserRoutes() {
	app.router.HandleFunc("/users", app.CreateUser).Methods("POST")
	app.router.HandleFunc("/users/{id}", app.GetUser).Methods("GET")
	app.router.HandleFunc("/users/{id}/profile", app.GetUserProfile).Methods("GET")
	app.router.HandleFunc("/users", app.UpdateUser).Methods("PATCH")
	app.router.HandleFunc("/users/{id}", app.DeleteUser).Methods("DELETE") // TODO: We will not be deleting data. We will only put an account in a deactivated state
	// app.router.HandleFunc("/users/{id}", app.GetAllUsers).Methods("GET")
}

func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser model.User
	reqBody, err := ioutil.ReadAll(r.Body) // Read the request body
	// TODO: Validate if the user already exist by checking the email ... here or on the side of postgres?
	if err != nil {
		log.Printf("App.CreateUser - error reading request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &newUser) // Fill newUser with the values coming from frontend
	if err != nil {
		log.Printf("App.CreateUser - error unmarshaling request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := app.store.UserProvider.CreateUser(&newUser)
	if err != nil {
		log.Printf("App.CreateUser - error creating user %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newUser.ID = id
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newUser)
	return
}

func (app *App) GetUser(w http.ResponseWriter, r *http.Request) {
	// Input
	userID := mux.Vars(r)["id"]
	// Validation
	// 1. Of a particular type
	//    i.e. check if its a string
	// 2. Particular format
	// 	  i.e. regex YYYY/MM/DD
	if userID == "" { // TODO: REGEX to validate other forms
		log.Printf("App.GetOneUser - empty user id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	usr, err := app.store.UserProvider.GetUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usr) // <- Sending the usr as a json {id: ..., first_name: ..., last_name ... , .. }
}

func (app *App) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Input
	userID := mux.Vars(r)["id"]
	// Validation
	// 1. Of a particular type
	//    i.e. check if its a string
	// 2. Particular format
	// 	  i.e. regex YYYY/MM/DD
	if userID == "" { // TODO: REGEX to validate other forms
		log.Printf("App.GetOneUser - empty user id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	usr, err := app.store.UserProvider.GetUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usr) // <- Sending the usr as a json {id: ..., first_name: ..., last_name ... , .. }
}

func (app *App) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Input - POST JSON
	// Validation
	// TODO
	var updatedUser model.User
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("App.UpdateUser - could not read r.Body with ioutil")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &updatedUser)
	if err != nil {
		log.Printf("App.UpdateUser - was unable to unmarshal changes")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// TODO: Validate that the updated user exists
	usr, err := app.store.UserProvider.UpdateUser(&updatedUser)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usr) // <- Sending the usr as a json {id: ..., first_name: ..., last_name ... , .. }
}

func (app *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	// TODO: More validation
	if userID == "" {
		log.Printf("App.RemoveUser - empty user id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.UserProvider.RemoveUser(userID)
	if err != nil {
		log.Printf("App.RemoveUser - error getting all users from provider %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}