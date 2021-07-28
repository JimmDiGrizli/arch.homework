package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

var db *gorm.DB

type User struct {
	Id        uint   `gorm:"primaryKey" json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

func main() {
	fmt.Println("-> Run server")
	fmt.Println("-> " + os.Getenv("PORT"))
	fmt.Println("-> " + os.Getenv("DATABASE_URI"))

	db, _ = gorm.Open(postgres.Open(os.Getenv("DATABASE_URI")), &gorm.Config{})

	router := mux.NewRouter()
	router.HandleFunc("/health", HealthCheck).Methods("GET")
	router.HandleFunc("/user", CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}", GetUser).Methods("GET")
	router.HandleFunc("/user/{id}", UpdateUser).Methods("PUT")
	router.HandleFunc("/user/{id}", DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func CreateUser(writer http.ResponseWriter, request *http.Request) {
	initHeaders(writer)

	var user User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	db.Create(&user)
	writer.WriteHeader(http.StatusOK)
}

func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	initHeaders(writer)

	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		log.Fatal(err.Error())
		return
	}

	db.Delete(&User{}, id)
	writer.WriteHeader(http.StatusNoContent)
}

func UpdateUser(writer http.ResponseWriter, request *http.Request) {
	initHeaders(writer)

	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var updatedUser User
	err = json.NewDecoder(request.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var user User
	result := db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	db.Model(&user).Updates(updatedUser)
	writer.WriteHeader(http.StatusOK)
}

func GetUser(writer http.ResponseWriter, request *http.Request) {
	initHeaders(writer)

	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var user User
	result := db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(writer).Encode(user)
	writer.WriteHeader(http.StatusOK)
}

func HealthCheck(writer http.ResponseWriter, _ *http.Request) {
	initHeaders(writer)
	fmt.Fprintf(writer, "{\"Status\":\"ok\"}")
	writer.WriteHeader(http.StatusOK)
}

func initHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}
