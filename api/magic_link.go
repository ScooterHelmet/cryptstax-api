package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"net/url"
	"strings"
	"time"
	"regexp"
	"strconv"
	"github.com/lib/pq"
)

const (
	verificationTTL		= time.Minute * 15
)

var (
	rxEmail    		= regexp.MustCompile("^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$")
	rxUUID			= regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")
)

var magicLinkTemplate = template.Must(template.ParseFiles("templates/magic-link.html"))
// MailSender provides a method to send mails.
type MailSender struct {
	Addr string
	Auth smtp.Auth
	From mail.Address
}
// Mail to send.
type Mail struct {
	To      mail.Address
	Subject string
	Body    string
}

func (s *server) handleMagicLink() http.HandlerFunc {
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

		// Handle VerificationCode
		var verificationCode string
		err := s.db.QueryRowContext(ctx, `
			INSERT INTO cryptstax.public.verification_codes (user_id) VALUES
				(
					(SELECT id FROM cryptstax.public.users WHERE email = $1)
				)
			RETURNING id`, input.Email,
			).Scan(
				&verificationCode,
			)

		if err != nil {
			magicLinkError := errors.New(err.(*pq.Error).Message)
			CheckError(magicLinkError)
			http.Error(w, err.(*pq.Error).Message, http.StatusInternalServerError)
			return
		}
		
		// Handle Magic Link
		magicLink := s.origin
		magicLink.Path = "/api/verify_redirect"
		q := magicLink.Query()
		q.Set("verification_code", verificationCode)
		q.Set("redirect_uri", input.RedirectURI)
		magicLink.RawQuery = q.Encode()

		var b bytes.Buffer
		data := map[string]interface{}{
			"MagicLink": magicLink.String(),
			"Minutes":   int(verificationTTL.Minutes()),
		}

		if err := magicLinkTemplate.Execute(&b, data); err != nil {
			http.Error(w, "could not execute magic link template: %v", http.StatusInternalServerError)
			return
		}

		if err := s.mailSender.send(Mail{
			To:      mail.Address{Address: input.Email},
			Subject: "Magic Link",
			Body:    b.String(),
		}); 
		
		err != nil {
			http.Error(w, "could not mail magic link: %v", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *server) verifyMagicRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	if err := s.db.PingContext(ctx); err != nil {
		CheckError(err)
	}

	q := r.URL.Query()
	verificationCode := q.Get("verification_code")
	redirectURI := q.Get("redirect_uri")

	errs := make(map[string]string)
	verificationCode = strings.TrimSpace(verificationCode)
	if verificationCode == "" {
		errs["verification_code"] = "Verification code required"
	} else if !rxUUID.MatchString(verificationCode) {
		errs["verification_code"] = "Invalid verification code"
	}

	var callback *url.URL
	var err error
	redirectURI = strings.TrimSpace(redirectURI)
	if redirectURI == "" {
		errs["redirect_uri"] = "Redirect URI required"
	} else if callback, err = url.ParseRequestURI(redirectURI); err != nil {
		errs["redirect_uri"] = "Invalid redirect URI"
	}
	if len(errs) != 0 {
		http.Error(w, "unable to process uri", http.StatusUnprocessableEntity)
		return
	}

	var userID string
	var createdAt time.Time
	if err := s.db.QueryRowContext(ctx, `
		DELETE FROM cryptstax.public.verification_codes WHERE id = $1
		RETURNING user_id, created_at`, 
		verificationCode,
		).Scan(
			&userID, 
			&createdAt,
			); 
		err == sql.ErrNoRows {
		http.Error(w, "Magic link not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "could not delete verification code: %v", http.StatusUnprocessableEntity)
		return
	}

	if createdAt.Add(verificationTTL).Before(time.Now()) {
		http.Error(w, "Link expired", http.StatusGone)
		return
	}

	var user User
	user.ID = userID
	userToken := CustomToken{}
	userToken.Token = CreateJWT(user)

	f := url.Values{}
	f.Set("token", userToken.Token)
	callback.Fragment = f.Encode()

	http.Redirect(w, r, callback.String(), http.StatusFound)
	}
}

/*
	Sending Mail to the SMTP Host server
*/
func (s *server) newMailSender(host string, port int, username, password string) *MailSender {
	return &MailSender{
		Addr: net.JoinHostPort(host, strconv.Itoa(port)),
		Auth: smtp.PlainAuth("", username, password, host),
		From: mail.Address{
			Name:    "cryptstax",
			Address: "noreply@" + s.origin.Hostname(),
		},
	}
}

func (s *MailSender) send(mail Mail) error {
	headers := map[string]string{
		"From":         s.From.String(),
		"To":           mail.To.String(),
		"Subject":      mail.Subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=utf-8",
	}
	msg := ""
	for k, v := range headers {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n"
	msg += mail.Body

	return smtp.SendMail(
		s.Addr,
		s.Auth,
		s.From.Address,
		[]string{mail.To.Address},
		[]byte(msg))
}