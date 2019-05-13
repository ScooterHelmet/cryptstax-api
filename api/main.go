package main

import (
	"context"
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/url"
	//"net/smtp"
)

type server struct {
	origin		*url.URL
	db			*sql.DB
	ctx context.Context
	router		*mux.Router
	//smtp	smtp.Auth	//MailTrap.io STMP_USERNAME and STMP_PASSWORD
}

func main() {
	//	@TODO: - connect blockchain networks - 
	//	implement Crypstax network endpoint
	//	implement Steemit network endpoint
	//	implement Ethereum network endpoint
	//	...

	server := &server{
		router: mux.NewRouter(),
	}
	server.ConnectDB()
	//server.ConnectCrypstaxNetwork()
	server.Start()
}
