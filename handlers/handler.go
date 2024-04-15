package handlers

import (
	"log"
	"net/http"

	"github.com/DivyanshuVerma98/goFileProcessing/constants"
	"github.com/gorilla/mux"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productType := params["product_type"]
	log.Println("Given productType -", productType)
	if len(productType) == 0 {
		SendResponse(w, "product_type key is required", http.StatusBadRequest, nil)
		return
	}
	switch productType {
	case constants.Motor:
		MotorService(w, r)
	default:
		SendResponse(w, "Invalid product_type selected", http.StatusBadRequest, nil)
	}
}
