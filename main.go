package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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

// Metadata Struct (Model)
type Metadata struct {
	LastModified     string `json:"lastmodified"`
	LastmodifiedDate string `json:"lastmodifieddate"`
	Size             string `json:"size"`
	Type             string `json:"type"`
}

// Init objects var as a slice Object struct
var objects []Object
var ecdsa []ECDSA
var config struct {
	port        int
	appURL      *url.URL
	databaseURL string
	jwtKey      []byte
	smtpAddr    string
	smtpAuth    smtp.Auth
}
var db *sql.DB

/*
func requireJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
			http.Error(w, "Content type of application/json required", http.StatusUnsupportedMediaType)
			return
		}
		next(w, r)
	}
}
*/

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

/**
* appURL will allow us to build the “magic link”.
* port in which the HTTP server will start.
* databaseURL is the CockroachDB address, I added /cryptstax_db to the previous address to indicate the database name.
* jwtKey used to sign JWTs.
* smtpAddr is a joint of SMTP_HOST + SMTP_PORT; we’ll use it to to send mails.
* smtpUsername and smtpPassword are the two required vars.
* smtpAuth is also used to send mails.
 */
func init() {
	config.port, _ = strconv.Atoi(env("PORT", "3000"))
	config.appURL, _ = url.Parse(env("APP_URL", "http://localhost:"+strconv.Itoa(config.port)+"/"))
	config.databaseURL = env("DATABASE_URL", "postgresql://root@127.0.0.1:26257/cryptstax_db?sslmode=disable")
	config.jwtKey = []byte(env("JWT_KEY", "super-duper-secret-key"))
	smtpHost := env("SMTP_HOST", "smtp.mailtrap.io")
	config.smtpAddr = net.JoinHostPort(smtpHost, env("SMTP_PORT", "25"))
	smtpUsername, ok := os.LookupEnv("SMTP_USERNAME")
	if !ok {
		log.Fatalln("could not find SMTP_USERNAME on environment variables")
	}
	smtpPassword, ok := os.LookupEnv("SMTP_PASSWORD")
	if !ok {
		log.Fatalln("could not find SMTP_PASSWORD on environment variables")
	}
	config.smtpAuth = smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
}

func env(key, fallbackValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallbackValue
	}
	return v
}

func mustEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%s required on environment variables", key))
	}
	return v
}

func intEnv(key string, fallbackValue int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallbackValue
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallbackValue
	}
	return i
}

func main() {

	// Init Route
	r := mux.NewRouter()

	var err error
	if db, err = sql.Open("postgres", config.databaseURL); err != nil {
		log.Fatalf("could not open database connection: %v\n", err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatalf("could not ping to database: %v\n", err)
	}

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

	// Passwordless Auth
	r.HandleFunc("/api/users", createUser).Methods("POST")
	r.HandleFunc("/api/passwordless/start", passwordlessStart).Methods("POST")
	r.HandleFunc("/api/passwordless/verify_redirect", passwordlessVerifyRedirect).Methods("GET")
	r.HandleFunc("/api/auth_user", guard(getAuthUser)).Methods("GET")
	//r.HandleFunc("POST", "/api/users", requireJSON(createUser))
	//r.HandleFunc("POST", "/api/passwordless/start", requireJSON(passwordlessStart))
	//r.HandleFunc("GET", "/api/passwordless/verify_redirect", passwordlessVerifyRedirect)
	//r.HandleFunc("GET", "/api/auth_user", guard(getAuthUser))

	log.Printf("starting server at %s 🚀\n", config.appURL)
	log.Fatalf("could not start server: %v\n", http.ListenAndServe(fmt.Sprintf(":%d", config.port), r))

	log.Fatal(http.ListenAndServe(":8000", r))
}
