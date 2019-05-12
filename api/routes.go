package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func (s *server) init() {

	// GET All
	s.router.HandleFunc("/api/channels", s.handleGetChannels()).Methods("GET")

	// GET {ID}
	s.router.HandleFunc("/api/channels/{id}", s.handleGetChannel()).Methods("GET")

	// POST
	s.router.HandleFunc("/api/channels", s.handleCreateChannel()).Methods("POST")

	// PUT {ID}
	s.router.HandleFunc("/api/channels/{id}", s.handleUpdateChannel()).Methods("PUT")

	// DELETE {ID}
	s.router.HandleFunc("/api/channels/{id}", s.handleDeleteChannel()).Methods("DELETE")

	// HEAD {ID}
	s.router.HandleFunc("/api/channels/{id}", s.handleCheckChannel()).Methods("HEAD")
}

// Start start api service
func (s *server) Start() {
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"})

	s.init()
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(s.router)))
	log.Printf("accepting connections on port: %d\n", ":8000")
}
