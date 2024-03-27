package handlers

import (
	"context"
	"log"
	"net/http"
)

// Need to define a custom type for context key
type customType string

const UserDataKey customType = "UserData"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authorization, exists := r.Header["Authorization"]
		log.Println("Auth token ", authorization)
		if !exists || len(authorization) < 1 {
			SendResponse(w, "Authorization key required", http.StatusUnauthorized, nil)
			return
		}
		userData, err := GetUserDetailsAPI(authorization[0])
		if err != nil {
			SendResponse(w, err.Error(), http.StatusBadRequest, nil)
			return
		}
		ctx := context.WithValue(r.Context(), UserDataKey, userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
