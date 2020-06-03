package utils

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"log"
	"mime/multipart"
)


// Validating image files to be of correct type

func CheckMime(imageFile multipart.File) (string, error) {

	mime, errMime := mimetype.DetectReader(imageFile)
	if errMime != nil {
		log.Println("App.CheckMime - Error handling mime: " + errMime.Error())
		return "", errMime
	}

	// MIME Reads part of the file, rewind to the start
	_, err := imageFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Println("App.CheckMime - Something went wrong with seeking back to the front")
		return "", err
	}

	if mime.String() == "image/png" {
		return ".png", nil
	} else if mime.String() == "image/jpeg" {
		return ".jpg", nil
	} else {
		return "", fmt.Errorf("invalid type")
	}
}


//// Generating uuid's when not through Postgres
//func UUIDGen() string {
//	return uuid.New()
//}