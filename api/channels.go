package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Struct (model)
type Workspace struct {
	Channel *Channel
}

type Channel struct {
	ID            string `json:"id"`
	Address       string `json:"name"`
	Created       string `json:"created"`
	Creator       string `json:"creator"`
	Is_archived   bool   `json:"is_archived"`
	Is_channel    bool   `json:"is_channel"`
	Is_general    bool   `json:"is_general"`
	Is_member     bool   `json:"is_member"`
	Is_mpim       bool   `json:"is_mpim"`
	Is_org_shared bool   `json:"is_org_shared"`
	Is_private    bool   `json:"is_private"`
	Is_shared     bool   `json:"is_shared"`
}

type Profile struct {
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalcode"`
	Country    string `json:"country"`
}

// Route Handlers
func (s *server) handleGetChannels() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(channels)
	}
}

func (s *server) handleGetChannel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r) // Get params
		// Loop through channels and find with id
		for _, item := range channels {
			if item.ID == params["id"] {
				json.NewEncoder(w).Encode(item)
				return
			}
		}
		json.NewEncoder(w).Encode(&Channel{})
	}
}

func (s *server) handleCreateChannel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var channel Channel
		_ = json.NewDecoder(r.Body).Decode(&channel)
		channel.ID = strconv.Itoa(rand.Intn(10000000)) //Mock ID - not safe/used in prod
		channels = append(channels, channel)
		json.NewEncoder(w).Encode(channel)
	}
}

func (s *server) handleUpdateChannel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		for index, item := range channels {
			if item.ID == params["id"] {
				channels = append(channels[:index], channels[index+1:]...)
				var channel Channel
				_ = json.NewDecoder(r.Body).Decode(&channel)
				channel.ID = params["id"]
				channels = append(channels, channel)
				json.NewEncoder(w).Encode(channel)
				return
			}
		}
		json.NewEncoder(w).Encode(channels)
	}
}

func (s *server) handleCheckChannel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		for _, item := range channels {
			if item.ID == params["id"] {
				return
			}
		}
		json.NewEncoder(w).Encode(&Channel{})
	}
}

func (s *server) handleDeleteChannel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		for index, item := range channels {
			if item.ID == params["id"] {
				channels = append(channels[:index], channels[index+1:]...)
				break
			}
		}
		json.NewEncoder(w).Encode(channels)
	}
}
