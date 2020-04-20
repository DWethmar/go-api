
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=api
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build: 
	$(GOBUILD) -v -o $(BINARY_NAME) ./cmd/api/
watch:
	modd
test: 
	$(GOTEST) -v ./...
benchmark:
	$(GOTEST) -bench=. ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
cert:
	mkdir -p cert
	openssl req -newkey rsa:2048 -nodes -keyout cert/server.key -x509 -days 365 -out cert/server.crt
run:
	$(GOBUILD) -v -o $(BINARY_NAME) ./cmd/api/ 
	./$(BINARY_NAME) -http 1
run-auth:
	$(GOBUILD) -v -o auth ./cmd/auth/ 
	./auth
deps:
	env GO111MODULE=on $(GOGET) github.com/cortesi/modd/cmd/modd

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
