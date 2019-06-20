package main

import (
	"context"
	"net/url"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
	"github.com/hako/branca"
)

type server struct {
	ctx			context.Context
	origin		*url.URL
	db			*sql.DB
	router		*mux.Router
	mailSender 	*MailSender
	codec      	*branca.Branca
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
	server.ConnectSMTP()
	server.Start()
}
