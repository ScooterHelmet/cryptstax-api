package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/dgryski/dgoogauth"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var emailConfig dgoogauth.OTPConfig

func (s *server) InitEmailAuthentication() {
	// Generate random secret
	secret := make([]byte, 10)
	_, err := rand.Read(secret)
	CheckError(err)

	secretBase32 := base32.StdEncoding.EncodeToString(secret)

	// The OTPConfig gets modified by otpc.Authenticate() to prevent passcode replay, etc.,
	// so allocate it once and reuse it for multiple calls.
	emailConfig = dgoogauth.OTPConfig{
		Secret:      secretBase32,
		WindowSize:  5,
		HotpCounter: 1,
		UTC:         true,
	}
}

func (s *server) HandleEmailCode(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		_ = json.NewDecoder(r.Body).Decode(&user)

		ctx := r.Context()
		if err := db.PingContext(ctx); err != nil {
			CheckError(err)
		}

		err := db.QueryRowContext(ctx,
			`SELECT
			email
			FROM cryptstax.public.users WHERE id = $1`, user.ID,
		).Scan(
			&user.Email,
		)

		if err != nil {
			CheckError(err)
			// Return error here if user wasn't found
			http.Error(w, "User not found", http.StatusBadRequest)
		}

		codeSent := EmailCode(user)

		if codeSent {
			w.Write([]byte("Email verification code sent"))
		} else {
			http.Error(w, "Error sending email", http.StatusInternalServerError)
		}
	}
}

func EmailCode(user User) bool {
	// Value parameter for ComputeCode must match window size in OTPConfig
	emailVerificationCode := dgoogauth.ComputeCode(emailConfig.Secret, int64(emailConfig.HotpCounter))

	from := mail.NewEmail("cryptstax", "cryptstax@no-reply.com")
	subject := "One time verification code for cryptstax"
	to := mail.NewEmail(user.Email, user.Email)
	plainTextConent := "cryptstax Verification Code"
	htmlContent := "Your cryptstax verification code is: " + strconv.Itoa(emailVerificationCode)
	message := mail.NewSingleEmail(from, subject, to, plainTextConent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	CheckError(err)

	// Client.Send returns 202 on success
	return response.StatusCode == 202
}

func (s *server) HandleVerifyEmail(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var code VerificationCode
		_ = json.NewDecoder(r.Body).Decode(&code)
		verified := verifyEmailCode(code.Code)

		if verified {
			// Find user in database to pass to CreateJWT
			ctx := r.Context()
			var user User

			if err := db.PingContext(ctx); err != nil {
				CheckError(err)
			}

			err := db.QueryRowContext(ctx,
				`SELECT id,
				email
				FROM cryptstax.public.users WHERE id = $1`, code.UserID,
			).Scan(
				&user.ID,
				&user.Email,
			)

			if err != nil {
				CheckError(err)
				// Return error here if user wasn't found
				http.Error(w, "User not found", http.StatusBadRequest)
			}

			userToken := CustomToken{}
			userToken.Token = CreateJWT(user)
			json.NewEncoder(w).Encode(userToken)
		} else {
			http.Error(w, "Email code is invalid", http.StatusBadRequest)
		}
	}
}

func verifyEmailCode(code string) bool {
	val, err := emailConfig.Authenticate(code)
	CheckError(err)

	if !val {
		return false
	}

	return true
}
