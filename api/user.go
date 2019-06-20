package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

// User represents an authenticated user or a resource owner.
type User struct {
	ID                 string `json:"id"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	PhoneNumber        string `json:"phone_number"`
	EmailVerification  bool   `json:"email_verification"`
	SMSVerification    bool   `json:"sms_verification"`
	GoogleVerification bool   `json:"google_verification"`
}

var users []User
var results []string

func (s *server) handleRegistration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		ctx := r.Context()
		var user User
		_ = json.NewDecoder(r.Body).Decode(&user)

		if err := s.db.PingContext(ctx); err != nil {
			CheckError(err)
		}

		hash := HandleCrypto(user.Password)

		  result, err := s.db.ExecContext(ctx,`
			INSERT INTO cryptstax.public.users (
				email,
				pass
			) VALUES ($1, $2, $3);`,
			user.Email,
			hash,
		  )

		if err != nil {
			registrationError := errors.New(err.(*pq.Error).Message)
			CheckError(registrationError)
			http.Error(w, err.(*pq.Error).Message, http.StatusInternalServerError)
			return
		}
		
		if result != nil {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (s *server) handleGetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		ctx := r.Context()
		var users []User

		if err := s.db.PingContext(ctx); err != nil {
			CheckError(err)
		}

		rows, err := s.db.QueryContext(ctx,
		`SELECT * FROM cryptstax.public.users;`)

		CheckError(err)

		for rows.Next() {
			var user User
			err = rows.Scan(
				&user.ID,
				&user.Email,
				&user.Password,
				&user.PhoneNumber,
				&user.EmailVerification,
				&user.SMSVerification,
				&user.GoogleVerification,)
			if err != nil {
				break
			}
			users = append(users, user)
		}

		// Check for errors during rows "Close".
		// This may be more important if multiple statements are executed
		// in a single batch and rows were written as well as read.
		if closeErr := rows.Close(); closeErr != nil {
			CheckError(err)
			http.Error(w, closeErr.Error(), http.StatusInternalServerError)
			return
		}

		// Check for row scan error.
		if err != nil {
			CheckError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check for errors during row iteration.
		if err = rows.Err(); err != nil {
			CheckError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}

func (s *server) handleDeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		ctx := r.Context()
		var err error

		if err := s.db.PingContext(ctx); err != nil {
			log.Fatal(err)
		  }

		sqlStatement := `DELETE FROM cryptstax.public.users
		WHERE id = $1;`
		
		_, err = s.db.Exec(sqlStatement, params["id"])

		CheckError(err)
	}
}