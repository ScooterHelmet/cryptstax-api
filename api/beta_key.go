package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"github.com/lib/pq"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func (s *server) handleBetaKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		

		var input struct {
			Email       string `json:"email"`
			RedirectURI string `json:"redirectURI"`
		}
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		
		errs := make(map[string]string)
		input.Email = strings.TrimSpace(input.Email)
		if input.Email == "" {
			errs["email"] = "Email required"
		} else if !rxEmail.MatchString(input.Email) {
			errs["email"] = "Invalid email"
		}
		input.RedirectURI = strings.TrimSpace(input.RedirectURI)
		if input.RedirectURI == "" {
			errs["redirectUri"] = "Redirect URI required"
		} else if _, err := url.ParseRequestURI(input.RedirectURI); err != nil {
			errs["redirectUri"] = "Invalid redirect URI"
		}
		if len(errs) != 0 {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}

		ctx := r.Context()

		if err := s.db.PingContext(ctx); err != nil {
			CheckError(err)
		}

		// Handle BetaKey
		var betaKey string
		err := s.db.QueryRowContext(ctx, `
			INSERT INTO cryptstax.public.beta_keys (user_id) VALUES
				(
					(SELECT id FROM cryptstax.public.users WHERE email = $1)
				)
			RETURNING id`, input.Email,
			).Scan(
				&betaKey,
			)

		if err != nil {
			betaKeyError := errors.New(err.(*pq.Error).Message)
			CheckError(betaKeyError)
			http.Error(w, err.(*pq.Error).Message, http.StatusInternalServerError)
			return
		}

		from := mail.NewEmail("cryptstax", "noreply@cryptstax.com")
		subject := "Beta Access Key to cryptstax"
		to := mail.NewEmail(input.Email, input.Email)
		plainTextContent := "Congratulations!"
		htmlContent :=  "<span>Your Beta Key is "+ betaKey + "</span><br /><em><span>This beta key never expires and can only be used once. Redeem <a href="+input.RedirectURI+" target='_blank' rel='noopener'>here</a></span>"
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
		response, err := client.Send(message)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(response.StatusCode)
			fmt.Println(response.Body)
			fmt.Println(response.Headers)
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func (s *server) verifyBetaKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		var input struct {
			BetaKey       string `json:"betakey"`
		}

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		errs := make(map[string]string)
		input.BetaKey = strings.TrimSpace(input.BetaKey)
		if input.BetaKey == "" {
			errs["betakey"] = "BetaKey required"
		} else if !rxEmail.MatchString(input.BetaKey) {
			errs["betakey"] = "Invalid BetaKey"
		}

		var createdAt time.Time
		if err := s.db.QueryRowContext(ctx, `
		DELETE FROM cryptstax.public.beta_keys WHERE id = $1
		RETURNING created_at`, 
		input.BetaKey,
		).Scan( 
			&createdAt,
			);

		err == sql.ErrNoRows {
		http.Error(w, "BetaKey not found", http.StatusNotFound)
		return
		} else if err != nil {
			http.Error(w, "could not delete beta key: %v", http.StatusUnprocessableEntity)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}