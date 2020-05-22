package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-api/model"
	"io/ioutil"
	"log"
	"net/http"
)

func (app *App) RegisterThemeRoutes() {
	app.router.HandleFunc("/themes")

	app.router.HandleFunc("/users", app.CreateUser).Methods("POST")
	app.router.HandleFunc("/users/{id}", app.GetUser).Methods("GET")
	app.router.HandleFunc("/users/{id}/profile", app.GetUserProfile).Methods("GET")
	app.router.HandleFunc("/users", app.UpdateUser).Methods("PATCH")

	app.router.HandleFunc("/users/{id}", app.DeleteUser).Methods("DELETE")

	app.router.HandleFunc("/users/{follower_id}/follows/{followed_id}", app.FollowUser).Methods("POST")
	app.router.HandleFunc("/users/{follower_id}/follows/{followed_id}", app.UnfollowUser).Methods("DELETE")

	app.router.HandleFunc("/users/{user_id}/interested/{project_id}", app.InterestedProject).Methods("POST")
	app.router.HandleFunc("/users/{user_id}/interested/{project_id}", app.UninterestedProject).Methods("DELETE")

	app.router.HandleFunc("/users/{user_id}/contributes/{project_id}", app.JoinProject).Methods("POST")
	app.router.HandleFunc("/users/{user_id}/contributes/{project_id}", app.QuitProject).Methods("DELETE")

}
