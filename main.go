package main

import (
	"go-api/routes"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	// get all env
	// create a json

	if port == "" {
		port = "8000"
		log.Println("Did not find env var PORT, defaulting to 8080")
	}
	dbUrl := os.Getenv("DATABASE_URL")
	// More environment variables to be added here
	//projectID := os.Getenv("YOUR_PROJECT_ID")
	//bucketName := os.Getenv("my-new-bucket")
	app := routes.NewApp(&routes.AppConfig{
		DbUrl:          dbUrl,
		//ProjectID:		projectID,
		//BucketName:		bucketName,
	})
	defer app.Close()
	err := app.Setup(port)
	log.Fatal(err)
}
