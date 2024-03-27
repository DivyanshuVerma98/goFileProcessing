package handlers

import (
	"log"
	"net/http"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	productType := r.FormValue("product_type")
	log.Println("Given productType -", productType)
	if len(productType) == 0 {
		SendResponse(w, "product_type key is required", http.StatusBadRequest, nil)
		return
	}
	switch productType {
	case "motor":
		MotorService(w, r)
	default:
		SendResponse(w, "Invalid product_type selected", http.StatusBadRequest, nil)
	}
}
