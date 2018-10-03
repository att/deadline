
DESTDIR ?= ~
LD_FLAGS="-X main.commit=$(shell git rev-parse HEAD) -X main.version=$(shell cat VERSION) -X main.builtby=$(shell whoami)@$(shell hostname)"

server_file	        ?=	"cmd/main.go"
server_output		?=	"deadline-server"

threshold		?=	"10"
packages = $(shell go list ./... | grep -v /vendor/)
packages_relative = $(shell go list ./... | grep -v /vendor/ | cut -d"/" -f4)
#check_complexity = $(shell scripts/check_complexity.sh -t $(threshold) -f $(pkg))
check_complexity = $(shell echo $(pkg))

all: init verify test build

init:
	dep ensure

verify: 
	@echo "Verifying..."
	@go vet $(packages)
	for pkg in $(packages_relative); do scripts/check_complexity.sh -t 10 -f "$$pkg"; done

build: test
	@echo "Building..."
	@go build -ldflags $(LD_FLAGS) -o ${server_output} ${server_file}

test: verify
	@echo "Testing..."
	@go test -cover -race $(shell go list ./... | grep -v /vendor/)

