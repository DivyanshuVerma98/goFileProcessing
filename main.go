package main

import (
	"fmt"
	"net/http"

	"github.com/DivyanshuVerma98/goFileProcessing/handlers"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Print("Starting ..... \n")

	router := mux.NewRouter()
	router.HandleFunc("/fms/upload_doc/", handlers.UploadFileHandler).Methods("POST")
	http.ListenAndServe(":3000", router)
}
