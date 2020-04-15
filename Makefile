
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
run:
	$(GOBUILD) -v -o $(BINARY_NAME) ./cmd/api/ 
	./$(BINARY_NAME)
deps:
	env GO111MODULE=on $(GOGET) github.com/cortesi/modd/cmd/modd

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
