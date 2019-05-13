package main

import (
	"encoding/json"
	"net/http"
	"log"

	//"github.com/gorilla/mux"
)

type Channel struct {
	Id            int	 `json:"id"`
	Address       string `json:"address"`
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

// Init channels as slice Channel struct
var channels []Channel

// Route Handlers
func (s *server) handleCreateChannel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctx := r.Context()

		var channel Channel
		_ = json.NewDecoder(r.Body).Decode(&channel)

		if err := s.db.PingContext(ctx); err != nil {
			log.Fatal(err)
		  }

		result, err := s.db.ExecContext(ctx,`
			INSERT INTO cryptstax_db.public.channels (
				address,
				created,
				creator,
				is_archived,
				is_channel,
				is_general,
				is_member,
				is_mpim,
				is_org_shared,
				is_private,
				is_shared
			) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`, 
			channel.Address,
			channel.Created,
			channel.Creator, 
			channel.Is_archived, 
			channel.Is_channel, 
			channel.Is_general, 
			channel.Is_member,
			channel.Is_mpim,
			channel.Is_org_shared,
			channel.Is_private,
			channel.Is_shared,
		)

		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(result)
	}
}

func (s *server) handleGetChannels() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctx := r.Context()

		if err := s.db.PingContext(ctx); err != nil {
			log.Fatal(err)
		  }

		rows, err := s.db.QueryContext(ctx,
		`SELECT * FROM cryptstax_db.public.channels;`)

		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var channel Channel
			err = rows.Scan(
				&channel.Id,
				&channel.Address,
				&channel.Created,
				&channel.Creator, 
				&channel.Is_archived, 
				&channel.Is_channel, 
				&channel.Is_general, 
				&channel.Is_member,
				&channel.Is_mpim,
				&channel.Is_org_shared,
				&channel.Is_private,
				&channel.Is_shared,)
			if err != nil {
				break
			}
			channels = append(channels, channel)
			w.WriteHeader(http.StatusOK)
		}
		// Check for errors during rows "Close".
		// This may be more important if multiple statements are executed
		// in a single batch and rows were written as well as read.
		if closeErr := rows.Close(); closeErr != nil {
			http.Error(w, closeErr.Error(), http.StatusInternalServerError)
			return
		}

		// Check for row scan error.
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check for errors during row iteration.
		if err = rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer rows.Close()
	}
}

func (s *server) handleGetChannel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctx := r.Context()

		if err := s.db.PingContext(ctx); err != nil {
			log.Fatal(err)
		  }
		//params := mux.Vars(r) // Get params
		// Loop through channels and find with id

		json.NewEncoder(w).Encode(&Channel{})
	}
}



func (s *server) handleUpdateChannel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		//params := mux.Vars(r)
		/*for index, item := range channels {
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
		json.NewEncoder(w).Encode(channels)*/
	}
}

func (s *server) handleCheckChannel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		//params := mux.Vars(r) // Get params
		// Loop through channels and find with id
		//for index, item := range channels {
			//if item.ID == params["id"] {
				//channels = append(channels[:index], channels[index+1:]...)
				//break
			//}
		//}
		//json.NewEncoder(w).Encode(channels)
	}
}

func (s *server) handleDeleteChannel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		//params := mux.Vars(r)
		//for index, item := range channels {
			//if item.ID == params["id"] {
				//channels = append(channels[:index], channels[index+1:]...)
				//break
			//}
		//}
		//json.NewEncoder(w).Encode(channels)
	}
}
