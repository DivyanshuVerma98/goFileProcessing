package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/DivyanshuVerma98/goFileProcessing/constants"
	"github.com/DivyanshuVerma98/goFileProcessing/structs"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing CORS middleware", r.Method)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, Origin, User-Agent, Cache-Control, Keep-Alive, X-Requested-With, If-Modified-Since, Token, Lang, Source")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func UserAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Calling UserAuthentication API to get user data using AuthToken
		time.Sleep(time.Second)
		userData := structs.UserData{
			ID:           007,
			FirstName:    "Divyasnhu",
			LastName:     "Verma",
			Username:     "DivyasnhuVerma98",
			BusinessRole: constants.MakerBusinessRole,
		}
		ctx := context.WithValue(r.Context(), constants.UserDataKey, &userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
