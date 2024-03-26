package handlers

import (
	"net/http"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	MotorService(w, r)
}
