
DESTDIR ?= ~

server_file	        ?=	"cmd/main.go"
server_output		?=	"deadline-server"

threshold		?=	"10"

all: init verify test build

init:
	dep ensure

verify: 
	@echo "Verifying..."
	@go vet $(shell go list ./... | grep -v /vendor/)
	#@sh scripts/check_complexity.sh -t ${threshold} -f ${server_file}


build: test
	@echo "Building..."
	@go build -o ${server_output} ${server_file}

test: verify
	@echo "Testing..."
	@go test $(shell go list ./... | grep -v /vendor/)

