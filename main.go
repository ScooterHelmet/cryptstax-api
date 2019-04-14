package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Object Struct (Model)
type Object struct {
	ID       string    `json:"id"`
	File     string    `json:"file"`
	FileName string    `json:"filename"`
	Metadata *Metadata `json:"metadata"`
}

// Metadata Struct
type Metadata struct {
	LastModified     string `json:"lastmodified"`
	LastmodifiedDate string `json:"lastmodifieddata"`
	Size             string `json:"size"`
	Type             string `json:"type"`
}

// Route Handlers
func getObjects(w http.ResponseWriter, r *http.Request) {

}
func getObject(w http.ResponseWriter, r *http.Request) {

}
func createObject(w http.ResponseWriter, r *http.Request) {

}
func updateObject(w http.ResponseWriter, r *http.Request) {

}
func checkObject(w http.ResponseWriter, r *http.Request) {

}
func deleteObject(w http.ResponseWriter, r *http.Request) {

}

func main() {
	// Init Route
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	r.HandleFunc("/api/objects", getObjects).Methods("GET")
	r.HandleFunc("/api/objects/{id}", getObject).Methods("GET")
	r.HandleFunc("/api/objects", createObject).Methods("POST")
	r.HandleFunc("/api/objects/{id}", updateObject).Methods("PUT")
	r.HandleFunc("/api/objects/{id}", checkObject).Methods("HEAD")
	r.HandleFunc("/api/objects/{id}", deleteObject).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
