package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (s *server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		ctx := r.Context()
		var user User
		_ = json.NewDecoder(r.Body).Decode(&user)
		providedPass := user.Password

		if err := s.db.PingContext(ctx); err != nil {
			CheckError(err)
		}

		err := s.db.QueryRowContext(ctx,
			`SELECT id,
			email,
			pass,
			phone_number,
			email_verification,
			sms_verification,
			google_verification	
			FROM cryptstax.public.users WHERE email = $1`, user.Email,
		).Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.PhoneNumber,
			&user.EmailVerification,
			&user.SMSVerification,
			&user.GoogleVerification,
		)
	
		if err != nil {
			// Email address sent by user doesn't exist in our database
			loginError := errors.New("user does not exist")
			CheckError(loginError)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		match, err := ComparePasswordAndHash(providedPass, user.Password)
		
		if (match) {
			// Check if user has 2FA enabled
			if user.EmailVerification || user.SMSVerification || user.GoogleVerification {
				if user.EmailVerification {
					// Send user the verification code by email
					emailSent := EmailCode(user)
					if emailSent {
						verificationData := UserVerification{
							UserID: user.ID,
							VerificationRoute: "/email/verify",
						}
						json.NewEncoder(w).Encode(verificationData)
						// Redirect to form to input code
					} else {
						http.Error(w, "Error sending email", http.StatusInternalServerError)
					}
				}

				 if user.SMSVerification {
				 	// Send user the verification code by sms
				 	messageSent := SendSMSCode(user)
				 	if messageSent {
				 		verificationData := UserVerification{
				 			UserID: user.ID,
				 			VerificationRoute: "/sms/verify",
				 		}
				 		json.NewEncoder(w).Encode(verificationData)
				 		// Redirect to form to input code
				 	} else {
				 		http.Error(w, "Error sending SMS", http.StatusInternalServerError)
				 	}
				 }

				if user.GoogleVerification {
					// Redirect to form to type in google authenticator code
					verificationData := UserVerification{
						UserID: user.ID,
						VerificationRoute: "/authenticator/verify",
					}
					json.NewEncoder(w).Encode(verificationData)
				}
			} else {
				// If no 2FA enabled, return authenticated jwt
				userToken := CustomToken{}
				userToken.Token = CreateJWT(user)
				json.NewEncoder(w).Encode(userToken)
			}
		} else {
			loginError := errors.New("invalid credentials")
			CheckError(loginError)
			http.Error(w, "invalid credentials", http.StatusForbidden)
			return
		}
	}
}