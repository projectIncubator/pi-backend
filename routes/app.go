package routes

import (
	"go-api/db"
	"go-api/db/postgres"
	"log"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type App struct {
	router *mux.Router
	store  *db.DataStore
	// TODO: Add Google Cloud Storage
	jwtMiddleware *jwtmiddleware.JWTMiddleware
}

type AppConfig struct {
	DbUrl string
	//ProjectID string TODO
	//BucketName string
}

func NewApp(config *AppConfig) *App {
	store, err := postgres.NewPostgresDataStore(config.DbUrl)
	// TODO: Create context and bucket instance?
	if err != nil {
		log.Println("Failed to connect to database")
		log.Fatal(err)
	}
	return &App{
		router:        mux.NewRouter().StrictSlash(true),
		store:         store,
		jwtMiddleware: InitAuthMiddleware(),
	}
}

func (app *App) Setup(port string) error {
	app.router.HandleFunc("/", app.index)
	app.RegisterRoutes()
	log.Println("App running at port:", port)
	handler := cors.AllowAll().Handler(app.router)
	return http.ListenAndServe(":"+port, handler)
}

func (app *App) RegisterRoutes() {
	// TODO: Will have i.e. below
	app.RegisterUserRoutes()
	// app.RegisterUserProfileRoutes()
	// app.RegisterProfileRoutes()
	app.RegisterProjectRoutes()
	app.RegisterThemeRoutes()
	app.RegisterGoogleCloudRoutes()
	app.RegisterDiscussionRoutes()
}

func (app *App) Close() {
	app.store.Close()
}

func (app *App) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"1","Hello World!",}`)) // TODO
}
