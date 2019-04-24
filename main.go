package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ECDSA Struct (Model)
type ECDSA struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

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
	LastmodifiedDate string `json:"lastmodifieddate"`
	Size             string `json:"size"`
	Type             string `json:"type"`
}

// Init objects var as a slice Object struct
var objects []Object
var ecdsa []ECDSA

// ECDSA Route Handlers
func getECDSAs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ecdsa)
}

// Route Handlers
func getObjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(objects)
}

func getObject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	// Loop through objects and find with id
	for _, item := range objects {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Object{})
}

func createObject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var object Object
	_ = json.NewDecoder(r.Body).Decode(&object)
	object.ID = strconv.Itoa(rand.Intn(10000000)) //Mock ID - not safe/used in prod
	objects = append(objects, object)
	json.NewEncoder(w).Encode(object)
}

func updateObject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range objects {
		if item.ID == params["id"] {
			objects = append(objects[:index], objects[index+1:]...)
			var object Object
			_ = json.NewDecoder(r.Body).Decode(&object)
			object.ID = params["id"]
			objects = append(objects, object)
			json.NewEncoder(w).Encode(object)
			return
		}
	}
	json.NewEncoder(w).Encode(objects)
}

func checkObject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range objects {
		if item.ID == params["id"] {
			return
		}
	}
	json.NewEncoder(w).Encode(&Object{})
}

func deleteObject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range objects {
		if item.ID == params["id"] {
			objects = append(objects[:index], objects[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(objects)
}

func main() {
	// Init Route
	r := mux.NewRouter()

	// Mock Data - @todo - implement Hyperledger Fabric endpoint
	objects = append(objects,
		Object{
			ID:       "1",
			File:     "asfbkjeafsabrnjekq7",
			FileName: "Hello_World1.txt",
			Metadata: &Metadata{
				LastModified:     "04/12/2019",
				LastmodifiedDate: "04/11/2019",
				Size:             "5647",
				Type:             "text"}})
	objects = append(objects,
		Object{ID: "2",
			File:     "asabrnkq7",
			FileName: "Hello_World2.txt",
			Metadata: &Metadata{
				LastModified:     "04/12/2019",
				LastmodifiedDate: "04/11/2019",
				Size:             "562347", Type: "text"}})
	objects = append(objects,
		Object{ID: "3",
			File:     "asbybyqqweqerw7",
			FileName: "Hello_World3.txt",
			Metadata: &Metadata{
				LastModified:     "04/12/2019",
				LastmodifiedDate: "04/11/2019",
				Size:             "562217", Type: "text"}})
	objects = append(objects,
		Object{ID: "4",
			File:     "asbybasdvsefvasdvrw7",
			FileName: "Hello_World4.txt",
			Metadata: &Metadata{
				LastModified:     "04/13/2019",
				LastmodifiedDate: "04/13/2019",
				Size:             "31227", Type: "text"}})

	// Objects -- Route Handlers / Endpoints
	r.HandleFunc("/api/objects", getObjects).Methods("GET")
	r.HandleFunc("/api/objects/{id}", getObject).Methods("GET")
	r.HandleFunc("/api/objects", createObject).Methods("POST")
	r.HandleFunc("/api/objects/{id}", updateObject).Methods("PUT")
	r.HandleFunc("/api/objects/{id}", checkObject).Methods("HEAD")
	r.HandleFunc("/api/objects/{id}", deleteObject).Methods("DELETE")

	// ECDSAs -- Route Handlers / Endpoints
	r.HandleFunc("/api/ecdsas", getECDSAs).Methods("GET")
	//r.HandleFunc("/api/ecdas/{id}", getECDSA).Methods("GET")
	//r.HandleFunc("/api/ecdsa", createECDSA).Methods("POST")
	//r.HandleFunc("/api/ecdsa/{id}", updateECDSA).Methods("PUT")
	//r.HandleFunc("/api/ecdsa/{id}", checkECDSA).Methods("HEAD")
	//r.HandleFunc("/api/ecdsa/{id}", deleteECDSA).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
