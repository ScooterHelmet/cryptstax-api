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

## Build API
```
cd ./cryptstax-api/api
go build
```

## Run API
```
./crypstax-api/api
```
* Navigate to http://localhost:8000 (404 response expected)
* Test with PostMan
