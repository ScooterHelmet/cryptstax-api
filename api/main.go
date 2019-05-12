package main

import (
	"context"
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/url"
	//"net/smtp"
	//"database/sql"
)

type server struct {
	origin		*url.URL
	db			*sql.DB
	ctx context.Context
	router		*mux.Router
	//smtp	smtp.Auth	//MailTrap.io STMP_USERNAME and STMP_PASSWORD
}

func main() {
	// @TODO: - connect blockchain network - implement Hyperledger Fabric endpoint

	/* TEST DATA
	channels = append(channels,
		Channel{
			Address:       "HTTP://127.0.0.1:7054",
			Created:       "May-10-2019 10:46:54 PM",
			Creator:       "Hyperledger Fabric",
			Is_archived:   false,
			Is_channel:    true,
			Is_general:    true,
			Is_member:     true,
			Is_mpim:       true,
			Is_org_shared: false,
			Is_private:    false,
			Is_shared:     true,
		})
	channels = append(channels,
		Channel{
			Address:       "HTTP://127.0.0.1:7545",
			Created:       "Jul-30-2015 03:26:13 PM",
			Creator:       "Ethereum",
			Is_archived:   false,
			Is_channel:    true,
			Is_general:    true,
			Is_member:     true,
			Is_mpim:       true,
			Is_org_shared: false,
			Is_private:    false,
			Is_shared:     true,
		})
		*/

	server := &server{
		router: mux.NewRouter(),
	}
	server.ConnectDB()
	server.Start()
}
