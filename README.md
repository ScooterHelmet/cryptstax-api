# Cryptstax API

## Prerequisites
* Golang Installed - https://golang.org/
```
// OSX
brew install go
```
* Gorilla Mux Installed - https://github.com/gorilla/mux
```
go get -u github.com/gorilla/mux
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
