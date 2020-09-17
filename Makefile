
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
SERVER_BINARY_NAME=server
SEEDER_BINARY_NAME=seeder

BIN=bin
SERVER_OUT=$(BIN)/$(SERVER_BINARY_NAME)
SEEDER_OUT=$(BIN)/$(SEEDER_BINARY_NAME)

all: test buildServer buildSeeder
buildServer: 
	$(GOBUILD) -v -o $(SERVER_OUT) ./cmd/server/
buildServerArm: 
	env GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -v -o $(SERVER_OUT) ./cmd/server/
buildSeeder:
	$(GOBUILD) -v -o $(SEEDER_OUT) ./cmd/seeder/
buildSeederArm:
	env GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -v -o $(SEEDER_OUT) ./cmd/seeder/
watch:
	modd
test: 
	$(GOTEST) -v ./...
benchmark:
	$(GOTEST) -bench=. ./...
clean: 
	$(GOCLEAN)
	rm -f $(SERVER_OUT)
	rm -f $(SEEDER_OUT)
run:
	$(GOBUILD) -v -o $(SERVER_BINARY_NAME) ./cmd/server/ 
	./$(SERVER_OUT) -port 8080
run-auth:
	$(GOBUILD) -v -o auth ./cmd/auth/ 
	./auth
deps:
	env GO111MODULE=on $(GOGET) github.com/cortesi/modd/cmd/modd
build-docker-arm:
	docker build -t go-api -f Dockerfile.raspbian .
