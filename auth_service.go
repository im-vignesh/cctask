package main

import (
	"./ccDB"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func authLoggingMiddleware(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

func main()  {
	r := mux.NewRouter()
	r.Use(authLoggingMiddleware)
	r.HandleFunc("/auth", authenticate).Methods("GET")
	err := http.ListenAndServe(":5081", r)
	if err != nil {
		log.Fatalf("Server Failed to Start: %v", err)
	}
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("Username")
	fmt.Println(username)
	db := ccDB.GetDBConnection()
	defer db.Close()
	var count int
	err := db.QueryRow("SELECT COUNT(username) FROM user WHERE username=?", username).Scan(&count)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if count == 1 {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
}

