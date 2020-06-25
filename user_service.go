package main

import (
	"./ccDB"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)


type User struct {
	Username string `json:"username"`
	DOB string `json:"dob"`
	Age int `json:"age"`
	Email string `json:"email"`
	PhoneNumber string `json:"phonenumber"`
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

func main() {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/user/profile", UserInfo).Methods("GET")
	r.HandleFunc("/microservice/name", MicroServiceName).Methods("GET")
	err := http.ListenAndServe(":5082", r)
	if err != nil {
		log.Fatalf("Server Failed to Start: %v", err)
	}
}

func UserInfo(w http.ResponseWriter, r *http.Request){
	username := r.Header.Get("Username")
	var user User
	db := ccDB.GetDBConnection()
	defer db.Close()
	var usernameSqlString,dob,email,phone_number sql.NullString
	err := db.QueryRow("SELECT username,dob,age, email, phone_number FROM user WHERE username=?", username).Scan(
		&usernameSqlString, &dob, &user.Age, &email, &phone_number)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	if usernameSqlString.Valid{
		user.Username = usernameSqlString.String
	}
	if dob.Valid{
		user.DOB = dob.String
	}
	if email.Valid{
		user.Email = email.String
	}
	if phone_number.Valid {
		user.PhoneNumber = phone_number.String
	}
	res, err := json.Marshal(user)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(res)
}

func MicroServiceName(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("user-microservice"))
}