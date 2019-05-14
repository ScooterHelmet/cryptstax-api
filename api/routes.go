package main

import (
	"log"
	"net/http"
	"os"
	"database/sql"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/joho/godotenv"
)

func (s *server) init() {
	// POST
	s.router.HandleFunc("/api/channels", s.handleCreateChannel()).Methods("POST")

	// GET All
	s.router.HandleFunc("/api/channels", s.handleGetChannels()).Methods("GET")

	// GET {ID}
	s.router.HandleFunc("/api/channels/{id}", s.handleGetChannelById()).Methods("GET")

	// PUT {ID}
	s.router.HandleFunc("/api/channels/{id}", s.handleUpdateChannelById()).Methods("PUT")

	// DELETE {ID}
	s.router.HandleFunc("/api/channels/{id}", s.handleDeleteChannel()).Methods("DELETE")
}

// Connect to external database service 
func (s *server) ConnectDB(){
	var (
		port         = intEnv("PORT", 8000)
		originStr    = env("ORIGIN", fmt.Sprintf("http://localhost:%d", port))
		dbURL        = env("DATABASE_URL", "postgresql://root@127.0.0.1:26257/?sslmode=disable")
		//secretKey    = env("SECRET_KEY", "supersecretkeyyoushouldnotcommit")
	)

	godotenv.Load()

	flag.IntVar(&port, "p", port, "Port ($PORT)")
	flag.StringVar(&originStr, "origin", originStr, "Origin ($ORIGIN)")
	flag.StringVar(&dbURL, "db", dbURL, "Database URL ($DATABASE_URL)")
	flag.Parse()

	var err error
	if s.origin, err = url.Parse(originStr); err != nil || !s.origin.IsAbs() {
		log.Fatalln("invalid origin")
		return
	}

	if i, err := strconv.Atoi(s.origin.Port()); err == nil {
		port = i
	}

	if s.db, err = sql.Open("postgres", dbURL); err != nil {
		log.Fatalf("could not open database connection: %v\n", err)
		return
	}

	//defer s.db.Close()

	if err = s.db.Ping(); err != nil {
		log.Fatalf("could not ping to database: %v\n", err)
		return
	}
}

// Start api service
func (s *server) Start() {
	handler := cors.Default().Handler(s.router)

	s.init()
	log.Fatal(http.ListenAndServe(":8000", handler))
	log.Printf("starting server at %s 🚀\n", s.origin.String())
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