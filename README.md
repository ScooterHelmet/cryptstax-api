# Cryptstax API

## Prerequisites
* Golang Installed - https://golang.org/
```
// OSX
brew install go
```
* GOPATH Dependencies
```
go get -u github.com/dgrijalva/jwt-go
go get -u github.com/dgryski/dgoogauth
go get -u github.com/gorilla/mux
go get -u github.com/hako/branca
go get -u github.com/joho/godotenv
go get -u github.com/lib/pq
go get -u github.com/rs/cors
go get -u github.com/sirupsen/logrus
go get -u github.com/sendgrid/sendgrid-go
go get -u github.com/sfreiberg/gotwilio
go get -u golang.org/x/crypto/argon2
go get -u gopkg.in/square/go-jose.v2
go get -u rsc.io/qr
```

## Add Project to GOPATH 
```
mkdir -p $GOPATH/src/github.com/[github_username]/cryptstax-api
```

## Build and Run API
```
cd ./cryptstax-api/api
go build
./api
```

## Testing with Mailtrap
```
cd ./cryptstax-api
touch .env
echo "SMTP_USERNAME=" >> .env
echo "SMTP_PASSWORD=" >> .env
```
* For local SMTP testing, be sure to paste in the username and password of your smtp.mailtrap.io credentials

* Navigate to http://localhost:8000 (404 response expected)
* Test with PostMan
