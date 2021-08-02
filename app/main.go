package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var db *gorm.DB
var random *rand.Rand

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

	s1 := rand.NewSource(time.Now().UnixNano())
	random = rand.New(s1)

	db, _ = gorm.Open(postgres.Open(os.Getenv("DATABASE_URI")), &gorm.Config{})

	router := mux.NewRouter()
	router.HandleFunc("/v1/health", HealthCheck).Methods("GET")
	router.Handle("/v1/metrics", promhttp.Handler()).Methods("GET")

	router.HandleFunc("/v1/user", CreateUser).Methods("POST")
	router.HandleFunc("/v1/user/{id}", GetUser).Methods("GET")
	router.HandleFunc("/v1/user/{id}", UpdateUser).Methods("PUT")
	router.HandleFunc("/v1/user/{id}", DeleteUser).Methods("DELETE")
	router.Use(metrics)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		metricsRequestCount.WithLabelValues(r.Method, r.RequestURI).Inc()

		next.ServeHTTP(w, r)

		metricsRequestLatency.
			WithLabelValues(r.Method, r.RequestURI, w.Header().Get("x-status-code")).
			Observe(time.Since(start).Seconds())
	})
}

func CreateUser(writer http.ResponseWriter, request *http.Request) {
	initHeaders(writer)

	if isError() {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("x-status-code", strconv.Itoa(http.StatusInternalServerError))
		return
	}

	var user User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		writer.Header().Set("x-status-code", strconv.Itoa(http.StatusBadRequest))
		return
	}

	db.Create(&user)
	fmt.Fprintf(writer, "{\"id\":"+strconv.Itoa(int(user.Id))+"}")
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("x-status-code", strconv.Itoa(http.StatusOK))
}

func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	initHeaders(writer)

	if isError() {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("x-status-code", strconv.Itoa(http.StatusInternalServerError))
		return
	}

	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		writer.Header().Set("x-status-code", strconv.Itoa(http.StatusBadRequest))
		log.Fatal(err.Error())
		return
	}

	db.Delete(&User{}, id)
	writer.WriteHeader(http.StatusNoContent)
	writer.Header().Set("x-status-code", strconv.Itoa(http.StatusNoContent))
}

func UpdateUser(writer http.ResponseWriter, request *http.Request) {
	initHeaders(writer)

	if isError() {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("x-status-code", strconv.Itoa(http.StatusInternalServerError))
		return
	}

	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		writer.Header().Set("x-status-code", strconv.Itoa(http.StatusBadRequest))
		return
	}

	var updatedUser User
	err = json.NewDecoder(request.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		writer.Header().Set("x-status-code", strconv.Itoa(http.StatusBadRequest))
		return
	}

	var user User
	result := db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		writer.WriteHeader(http.StatusNotFound)
		writer.Header().Set("x-status-code", strconv.Itoa(http.StatusNotFound))
		return
	}

	db.Model(&user).Updates(updatedUser)
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("x-status-code", strconv.Itoa(http.StatusOK))
}

func GetUser(writer http.ResponseWriter, request *http.Request) {
	initHeaders(writer)

	if isError() {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Header().Set("x-status-code", strconv.Itoa(http.StatusInternalServerError))
		return
	}

	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		writer.Header().Set("x-status-code", strconv.Itoa(http.StatusBadRequest))
		return
	}

	var user User
	result := db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		writer.WriteHeader(http.StatusNotFound)
		writer.Header().Set("x-status-code", strconv.Itoa(http.StatusNotFound))
		return
	}

	json.NewEncoder(writer).Encode(user)
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("x-status-code", strconv.Itoa(http.StatusOK))
}

func HealthCheck(writer http.ResponseWriter, _ *http.Request) {
	initHeaders(writer)
	fmt.Fprintf(writer, "{\"Status\":\"ok\"}")
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("x-status-code", strconv.Itoa(http.StatusOK))
}

func initHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func isError() bool {
	if random.Intn(1000) == 42 {
		return true
	} else {
		return false
	}
}
