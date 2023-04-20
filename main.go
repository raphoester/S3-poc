package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	router := mux.NewRouter()
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {}).Methods("GET")
	router.HandleFunc("/upload/{name}", Upload).Methods("GET")
	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Failed starting server | %s", err.Error())
	}
}

func Upload(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Bytes       []byte `json:"bytes"`
		ContentType string `json:"content_type"`
		FileName    string `json:"file_name"`
	}

	payload := input{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, fmt.Sprintf("failed decoding input body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	s := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				SharedConfigState: session.SharedConfigEnable,
			},
		),
	)

	s3Uploader := s3manager.NewUploader(s)

	if _, err := s3Uploader.Upload(
		&s3manager.UploadInput{
			Bucket:      aws.String("2im-demo-bucket-test"),
			Key:         aws.String(payload.FileName),
			Body:        bytes.NewReader(payload.Bytes),
			ContentType: aws.String(payload.ContentType),
		},
	); err != nil {
		http.Error(w, fmt.Sprintf("failed uploading image: %s", err.Error()), http.StatusInternalServerError)
		return
	}

}
