package routes

import (
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"go-api/model"
	"io/ioutil"
	"log"
	"net/http"
)

func (app *App) RegisterUserRoutes() {

	// Private APIs
	app.router.Handle("/api/private", negroni.New(
		negroni.HandlerFunc(app.jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			message := "Hello from a private endpoint! You need to be authenticated to see this."
			json.NewEncoder(w).Encode(message)
			w.WriteHeader(http.StatusOK)
		}))))
	app.router.Handle("/users", negroni.New(
		negroni.HandlerFunc(app.jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(app.CreateUser)))).Methods("POST")
	app.router.Handle("/users", app.middleware(http.HandlerFunc(app.UpdateProject),USER)).Methods("PATCH")
	app.router.Handle("/users", app.middleware(http.HandlerFunc(app.DeleteUser),USER)).Methods("DELETE")

	app.router.Handle("/follows/{followed_id}", negroni.New(
		negroni.HandlerFunc(app.jwtMiddleware.HandlerWithNext),
		negroni.Wrap(app.middleware(http.HandlerFunc(app.FollowUser),USER)))).Methods("POST")
	app.router.Handle("/follows/{followed_id}", negroni.New(
		negroni.HandlerFunc(app.jwtMiddleware.HandlerWithNext),
		negroni.Wrap(app.middleware(http.HandlerFunc(app.UnfollowUser),USER)))).Methods("DELETE")

	app.router.Handle("/interested/{project_id}", app.middleware(http.HandlerFunc(app.InterestedProject),USER)).Methods("POST")
	app.router.Handle("/interested/{project_id}", app.middleware(http.HandlerFunc(app.UninterestedProject),USER)).Methods("DELETE")

	app.router.Handle("/interested/{theme_name}", app.middleware(http.HandlerFunc(app.InterestedTheme),USER)).Methods("POST")
	app.router.Handle("/interested/{theme_name}", app.middleware(http.HandlerFunc(app.UninterestedTheme),USER)).Methods("DELETE")

	app.router.Handle("/contributes/{project_id}", app.middleware(http.HandlerFunc(app.JoinProject),USER)).Methods("POST")
	app.router.Handle("/contributes/{project_id}", app.middleware(http.HandlerFunc(app.QuitProject),USER)).Methods("DELETE")

	// Public APIs

	app.router.Handle("/users/{id}", http.HandlerFunc(app.GetUser)).Methods("GET")
	app.router.HandleFunc("/users/{id}/profile", app.GetUserProfile).Methods("GET")
	app.router.HandleFunc("/users/{id}/followers", app.GetUserFollowers).Methods("GET")
	app.router.HandleFunc("/users/{id}/follows", app.GetUserFollows).Methods("GET")
}

// Private APIs

func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {

	var response model.SignInResponse
	response.IsNewUser = false

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

	id, err := app.store.UserProvider.LoginUser(&newUser)
	if err != nil {

		id, err = app.store.UserProvider.CreateUser(&newUser)
		if err != nil {
			log.Printf("App.CreateUser - error creating user %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response.IsNewUser = true

	}

	response.UserID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
func (app *App) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var updatedUser model.UserProfile

	userID := r.Context().Value("user_id").(AuthWraper).id

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

	usr, err := app.store.UserProvider.UpdateUser(userID, &updatedUser)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usr) // <- Sending the usr as a json {id: ..., first_name: ..., last_name ... , .. }
}
func (app *App) DeleteUser(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("user_id").(AuthWraper).id

	err := app.store.UserProvider.RemoveUser(userID)
	if err != nil {
		log.Printf("App.RemoveUser - error removing the user %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) FollowUser(w http.ResponseWriter, r *http.Request) {

	followerID := r.Context().Value("user_id").(AuthWraper).id

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

	followerID := r.Context().Value("user_id").(AuthWraper).id

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

	userID := r.Context().Value("user_id").(AuthWraper).id

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

	userID := r.Context().Value("user_id").(AuthWraper).id

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

	userID := r.Context().Value("user_id").(AuthWraper).id

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

	userID := r.Context().Value("user_id").(AuthWraper).id


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

	userID := r.Context().Value("user_id").(AuthWraper).id

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

	userID := r.Context().Value("user_id").(AuthWraper).id

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
