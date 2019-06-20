package main

import (
	"log"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"net/url"
	"strconv"
	"database/sql"
	"github.com/rs/cors"
	"github.com/joho/godotenv"
	"gopkg.in/square/go-jose.v2/jwt"
)

type MiddlewareFunc func(http.Handler) http.Handler

func unrestrictedMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
    })
}

func (s *server) unrestrictedAPI() {
	unrestrictedRequests := s.router.PathPrefix("/api").Subrouter()
	unrestrictedRequests.Use(unrestrictedMiddleware)
	// REGISTRATION
	unrestrictedRequests.HandleFunc("/registration", s.handleRegistration()).Methods("POST")
	// LOGIN
	unrestrictedRequests.HandleFunc("/login", s.handleLogin()).Methods("POST")
	// MAGIC LINK
	unrestrictedRequests.HandleFunc("/magic_link", s.handleMagicLink()).Methods("POST")
	// AUTH MAGIC USER
	unrestrictedRequests.HandleFunc("/verify_redirect", s.verifyMagicRedirect()).Methods("GET")
	// Verify BETA KEY
	unrestrictedRequests.HandleFunc("/verify_betakey", s.verifyBetaKey()).Methods("POST")
	// VERIFY 2FA CODE FROM GOOGLE AUTHENTICATOR
	unrestrictedRequests.HandleFunc("/authenticator/verify", s.HandleVerifyGoogleAuthenticator(s.db)).Methods("POST")
	// VERIFY 2FA CODE SENT BY SMS
	unrestrictedRequests.HandleFunc("/sms/verify", s.HandleVerifySMS(s.db)).Methods("POST")
	// VERIFY 2FA CODE SENT BY EMAIL
	unrestrictedRequests.HandleFunc("/email/verify", s.HandleVerifyEmail(s.db)).Methods("POST")
}

func restrictedMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// HandleVerifyJWT takes an HTTP request and verifies if the correct header is sent with a valid JWT
		requestToken := r.Header.Get("Authorization")
		if requestToken != "" {
			splitToken := strings.Split(requestToken, " ")
			requestToken = splitToken[1]
			parsedJWT, err := jwt.ParseSigned(requestToken)
			CheckError(err)
	
			claim := CustomClaim{}
			err = parsedJWT.Claims(&PrivateRSAKey.PublicKey, &claim)
			CheckError(err)

			if err != nil {
				log.Println(r.RequestURI)
				http.Error(w, "Unauthorized Access", 401 )
				return
			}
		} else {
			log.Println(r.RequestURI)
			http.Error(w, "Unauthorized Access", 401 )
			return
		}
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (s *server) restrictedAPI() {
	restrictedRequests := s.router.PathPrefix("/api").Subrouter()
	restrictedRequests.Use(restrictedMiddleware)
	// POST BetaKey
	restrictedRequests.HandleFunc("/beta_key", s.handleBetaKey()).Methods("POST")
	// GET All
	restrictedRequests.HandleFunc("/users", s.handleGetUsers()).Methods("GET")
	// DELETE {ID}
	restrictedRequests.HandleFunc("/users/{id}", s.handleDeleteUser()).Methods("DELETE")
	// UPDATE 2FA SETTINGS
	restrictedRequests.HandleFunc("/users/2fa", s.handleUpdate2FA()).Methods("POST")
	// GENERATE QR CODE TO SCAN INTO GOOGLE AUTHENTICATOR
	restrictedRequests.HandleFunc("/authenticator/create", s.HandleGenerateQRCode(s.db)).Methods("GET")
	// SEND USER 2FA CODE BY SMS
	restrictedRequests.HandleFunc("/sms_code/send", s.HandleSendSMSCode(s.db)).Methods("POST")
	// SEND USER 2FA CODE BY EMAIL
	restrictedRequests.HandleFunc("/email_code/send", s.HandleEmailCode(s.db)).Methods("POST")
}

// Connect to external database service 
func (s *server) ConnectDB(){
	
	var (
		port         = intEnv("PORT", 8000)
		originStr    = env("ORIGIN", fmt.Sprintf("http://localhost:%d", port))
		dbURL        = env("DATABASE_URL", "postgresql://root@127.0.0.1:26257/?sslmode=disable")
	)

	godotenv.Load()

	flag.IntVar(&port, "p", port, "Port ($PORT)")
	flag.StringVar(&originStr, "origin", originStr, "Origin ($ORIGIN)")
	flag.StringVar(&dbURL, "db", dbURL, "Database URL ($DATABASE_URL)")
	flag.Parse()

	var err error
	if s.origin, err = url.Parse(originStr); err != nil || !s.origin.IsAbs() {
		log.Fatalln("invalid origin")
		return
	}

	if i, err := strconv.Atoi(s.origin.Port()); err == nil {
		port = i
	}

	if s.db, err = sql.Open("postgres", dbURL); err != nil {
		log.Fatalf("could not open database connection: %v\n", err)
		return
	}

	if err = s.db.Ping(); err != nil {
		log.Fatalf("could not ping to database: %v\n", err)
		return
	}
}

func (s *server) ConnectSMTP(){
	godotenv.Load()

	var (
		clientPort		= intEnv("CLIENT_PORT", 3000)
		clientURL		= env("CLIENT_URL", fmt.Sprintf("http://localhost:%d", clientPort))
		smtpHost		= env("SMTP_HOST", "smtp.mailtrap.io")
		smtpPort		= intEnv("SMTP_PORT", 25)
		smtpUsername	= mustEnv("SMTP_USERNAME")
		smtpPassword	= mustEnv("SMTP_PASSWORD")
	)

	flag.StringVar(&smtpHost, "smtp.host", smtpHost, "SMTP Host ($SMTP_HOST)")
	flag.IntVar(&smtpPort, "smtp.port", smtpPort, "SMTP Port ($SMTP_PORT)")
	flag.IntVar(&clientPort, "client.port", clientPort, "Client Port ($CLIENT_PORT)")
	flag.StringVar(&clientURL, "client.url", clientURL, "Client URL ($CLIENT_URL)")
	flag.Parse()

	s.mailSender = s.newMailSender(smtpHost, smtpPort, smtpUsername, smtpPassword)
}

// Start start api service
func (s *server) Start() {

	handler := cors.Default().Handler(s.router)
	
	s.unrestrictedAPI()
	s.restrictedAPI()
	s.InitJWT()
	s.InitLogrusPipeline()
	s.InitGoogleAuthenticator()
	s.InitTwilio()
	s.InitEmailAuthentication()

	log.Fatal(http.ListenAndServe(":8000", handler))
	log.Printf("accepting connections on port: %s\n", ":8000")
}