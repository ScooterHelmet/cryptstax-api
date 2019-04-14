# Netstack API

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
mkdir -p $GOPATH/src/github.com/[github_username]/netstack-api
```

## Build Project
```
cd ./netstack-api
go build
```

## Run Project
```
./netstack-api
```
* Navigate to http://localhost:8000 (404 response expected)
