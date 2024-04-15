package main

import (
	"fmt"
	"net/http"

	"github.com/DivyanshuVerma98/goFileProcessing/handlers"
	"github.com/DivyanshuVerma98/goFileProcessing/middleware"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Print("Starting ..... \n")

	router := mux.NewRouter()
	router.HandleFunc("/fms/upload_doc/{product_type}/", handlers.UploadFileHandler).Methods("POST")
	router.Use(middleware.UserAuthentication)
	http.ListenAndServe(":3000", middleware.CORSMiddleware(router))
}
