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
	app.router.HandleFunc("/themes", app.CreateTheme).Methods("POST")
	app.router.HandleFunc("/themes/{theme_name}", app.GetTheme).Methods("GET")
	app.router.HandleFunc("/themes/{theme_name}", app.UpdateTheme).Methods("PATCH")

	app.router.HandleFunc("/themes/{theme_name}", app.GetProjectsWithTheme).Methods("GET")

	app.router.HandleFunc("/themes/{theme_name}", app.DeleteTheme).Methods("DELETE")
}
