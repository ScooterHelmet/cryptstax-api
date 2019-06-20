package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/dgryski/dgoogauth"
	"rsc.io/qr"
)

var issuer = "blackbird"
var qrFilename = filepath.FromSlash("tmp/qr.png")
var totpConfig dgoogauth.OTPConfig

func (s *server) InitGoogleAuthenticator() {
	// Generate random secret
	secret := make([]byte, 10)
	_, err := rand.Read(secret)
	CheckError(err)

	secretBase32 := base32.StdEncoding.EncodeToString(secret)

	// The OTPConfig gets modified by otpc.Authenticate() to prevent passcode replay, etc.,
	// so allocate it once and reuse it for multiple calls.
	totpConfig = dgoogauth.OTPConfig{
		Secret:      secretBase32,
		WindowSize:  5,
		HotpCounter: 0,
		UTC:         true,
	}

	// Create directory to store temp files
	if !DirectoryExists("tmp") {
		err = os.Mkdir("tmp", 0777)
		CheckError(err)
	}
}

func (s *server) HandleGenerateQRCode(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		_ = json.NewDecoder(r.Body).Decode(&user)
		ctx := r.Context()

		if err := db.PingContext(ctx); err != nil {
			CheckError(err)
		}

		err := db.QueryRowContext(ctx,
			`SELECT email
				FROM cryptstax.public.users WHERE id = $1`, user.ID,
		).Scan(
			&user.Email,
		)

		if err != nil {
			CheckError(err)
			// Return error here if user wasn't found
			http.Error(w, "User not found", http.StatusBadRequest)
		}
		qrFile := s.generateQRCode(user.Email)
		http.ServeFile(w, r, qrFile)
	}
}

func (s *server) generateQRCode(email string) string {
	URL, err := url.Parse("otpauth://totp")
	CheckError(err)

	URL.Path += "/" + url.PathEscape(issuer) + ":" + url.PathEscape(email)

	params := url.Values{}
	params.Add("secret", totpConfig.Secret)
	params.Add("issuer", issuer)

	URL.RawQuery = params.Encode()
	code, err := qr.Encode(URL.String(), qr.Q)
	CheckError(err)

	b := code.PNG()
	err = ioutil.WriteFile(qrFilename, b, 0777)
	CheckError(err)

	return qrFilename
}

func (s *server) HandleVerifyGoogleAuthenticator(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var code VerificationCode
		_ = json.NewDecoder(r.Body).Decode(&code)
		validCode := s.verify2FACode(code.Code)

		if validCode {
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
			http.Error(w, "Google authenticator code is invalid", http.StatusBadRequest)
		}
	}
}

func (s *server) verify2FACode(code string) bool {
	val, err := totpConfig.Authenticate(code)
	CheckError(err)

	if !val {
		return false
	}

	return true
}
