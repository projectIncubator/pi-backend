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

	// Private APIs

	app.router.HandleFunc("/users", app.CreateUser).Methods("POST")
	app.router.HandleFunc("/users/{id_token}", app.UpdateUser).Methods("PATCH")
	app.router.HandleFunc("/users/{id_token}", app.DeleteUser).Methods("DELETE")

	app.router.HandleFunc("/users/{follower_token}/follows/{followed_id}", app.FollowUser).Methods("POST")
	app.router.HandleFunc("/users/{follower_token}/follows/{followed_id}", app.UnfollowUser).Methods("DELETE")

	app.router.HandleFunc("/users/{user_token}/interested/{project_id}", app.InterestedProject).Methods("POST")
	app.router.HandleFunc("/users/{user_token}/interested/{project_id}", app.UninterestedProject).Methods("DELETE")

	app.router.HandleFunc("/users/{user_token}/interested/{theme_name}", app.InterestedTheme).Methods("POST")
	app.router.HandleFunc("/users/{user_token}/interested/{theme_name}", app.UninterestedTheme).Methods("DELETE")

	app.router.HandleFunc("/users/{user_token}/contributes/{project_id}", app.JoinProject).Methods("POST")
	app.router.HandleFunc("/users/{user_token}/contributes/{project_id}", app.QuitProject).Methods("DELETE")

	// Public APIs

	app.router.HandleFunc("/users/{id}", app.GetUser).Methods("GET")
	app.router.HandleFunc("/users/{id}/profile", app.GetUserProfile).Methods("GET")
	app.router.HandleFunc("/users/{id}/followers", app.GetUserFollowers).Methods("GET")
	app.router.HandleFunc("/users/{id}/follows", app.GetUserFollows).Methods("GET")
}

// Private APIs

func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser model.IDUser
	reqBody, err := ioutil.ReadAll(r.Body) // Read the request body
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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(id)
}
func (app *App) UpdateUser(w http.ResponseWriter, r *http.Request) {

	userToken := mux.Vars(r)["id_token"]

	var updatedUser model.UserProfile
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

	usr, err := app.store.UserProvider.UpdateUser(userToken, &updatedUser)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usr) // <- Sending the usr as a json {id: ..., first_name: ..., last_name ... , .. }
}
func (app *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id_token"]
	if userID == "" {
		log.Printf("App.RemoveUser - empty user id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.UserProvider.RemoveUser(userID)
	if err != nil {
		log.Printf("App.RemoveUser - error removing the user %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) FollowUser(w http.ResponseWriter, r *http.Request) {

	followerID := mux.Vars(r)["follower_token"]
	followedID := mux.Vars(r)["followed_id"]

	if followerID == "" || followedID == "" {
		log.Printf("App.FollowUser - error reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := app.store.UserProvider.FollowUser(followerID, followedID)
	if err != nil {
		log.Printf("%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}
func (app *App) UnfollowUser(w http.ResponseWriter, r *http.Request) {

	followerID := mux.Vars(r)["follower_token"]
	followedID := mux.Vars(r)["followed_id"]

	if followerID == "" || followedID == "" {
		log.Printf("App.FollowUser - error reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.UserProvider.UnfollowUser(followerID, followedID)
	if err != nil {
		log.Printf("App.FollowUser - error creating user %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}

func (app *App) InterestedProject(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_token"]
	projectID := mux.Vars(r)["project_id"]

	if userID == "" || projectID == "" {
		log.Printf("App.FollowUser - error reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.UserProvider.InterestedProject(userID, projectID)
	if err != nil {
		log.Printf("App.FollowUser - error creating user %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}
func (app *App) UninterestedProject(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_token"]
	projectID := mux.Vars(r)["project_id"]

	if userID == "" || projectID == "" {
		log.Printf("App.FollowUser - error reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.UserProvider.UninterestedProject(userID, projectID)
	if err != nil {
		log.Printf("App.FollowUser - error creating user %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}

func (app *App) InterestedTheme(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_token"]
	themeName := mux.Vars(r)["theme_name"]

	if userID == "" || themeName == "" {
		log.Printf("App.FollowUser - error reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.UserProvider.InterestedTheme(userID, themeName)
	if err != nil {
		log.Printf("App.FollowUser - error creating user %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}
func (app *App) UninterestedTheme(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_token"]
	projectID := mux.Vars(r)["project_id"]

	if userID == "" || projectID == "" {
		log.Printf("App.FollowUser - error reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.UserProvider.UninterestedProject(userID, projectID)
	if err != nil {
		log.Printf("App.FollowUser - error creating user %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}

func (app *App) JoinProject(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_token"]
	projectID := mux.Vars(r)["project_id"]

	if userID == "" || projectID == "" {
		log.Printf("App.FollowUser - error reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.UserProvider.JoinProject(userID, projectID)
	if err != nil {
		log.Printf("App.FollowUser - error creating user %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}
func (app *App) QuitProject(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_token"]
	projectID := mux.Vars(r)["project_id"]

	if userID == "" || projectID == "" {
		log.Printf("App.FollowUser - error reading request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := app.store.UserProvider.QuitProject(userID, projectID)
	if err != nil {
		log.Printf("App.FollowUser - error creating user %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}

// Public APIs

func (app *App) GetUser(w http.ResponseWriter, r *http.Request) {
	// Input
	userID := mux.Vars(r)["id"]

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
	if userID == "" { // TODO: REGEX to validate other forms
		log.Printf("App.GetOneUser - empty user id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	usr, err := app.store.UserProvider.GetUserProfile(userID)
	if err != nil {
		log.Printf("%v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usr) // <- Sending the usr as a json {id: ..., first_name: ..., last_name ... , .. }
}
func (app *App) GetUserFollowers(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	if userID == "" { // TODO: REGEX to validate other forms
		log.Printf("App.GetOneUser - empty user id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	followers, err := app.store.UserProvider.GetUserFollowers(userID)
	if err != nil {
		log.Printf("%v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(followers) // <- Sending the usr as a json {id: ..., first_name: ..., last_name ... , .. }
}
func (app *App) GetUserFollows(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	if userID == "" { // TODO: REGEX to validate other forms
		log.Printf("App.GetOneUser - empty user id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	follows, err := app.store.UserProvider.GetUserFollows(userID)
	if err != nil {
		log.Printf("%v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(follows) // <- Sending the usr as a json {id: ..., first_name: ..., last_name ... , .. }
}
