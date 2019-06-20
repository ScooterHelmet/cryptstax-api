package main

import (
	"encoding/json"
	"net/http"
)

type VerificationCode struct {
	UserID string `json:"id"`
	Code   string `json:"code"`
}
// If a user has 2FA enabled this will hold those settings to be returned to the front end
type UserVerification struct {
	UserID string `json:"user_id"`
	VerificationRoute string `json:"verification_type"`
}

func (s *server) handleUpdate2FA() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		_ = json.NewDecoder(r.Body).Decode(&user)

		ctx := r.Context()
		err := s.db.PingContext(ctx)
		CheckError(err)

		sqlStatement := `UPDATE cryptstax.public.users
		SET phone_number = $2,
		email_verification = $3,
		sms_verification = $4,
		google_verification = $5
		WHERE email = $1;`

		_, err = s.db.Exec(
			sqlStatement,
			user.Email,
			user.PhoneNumber,
			user.EmailVerification,
			user.SMSVerification,
			user.GoogleVerification,
		)

		CheckError(err)

		if err != nil {
			http.Error(w, "Error updating 2FA settings", http.StatusInternalServerError)
		} else {
			w.Write([]byte("2FA settings updated successfully"))
		}
	}
}