
DESTDIR ?= ~
LD_FLAGS="-X main.commit=$(shell git rev-parse HEAD) -X main.version=$(shell cat VERSION) -X main.builtby=$(shell whoami)@$(shell hostname)"

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
	@go build -ldflags $(LD_FLAGS) -o ${server_output} ${server_file}

test: verify
	@echo "Testing..."
	@go test $(shell go list ./... | grep -v /vendor/)

