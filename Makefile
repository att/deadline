
DESTDIR ?= ~

server_file	        ?=	"cmd/deadline-server.go"
server_output		?=	"deadline-server"

threshold		?=	"10"

all: verify test build

verify: 
	@go vet ./...
	#@sh scripts/check_complexity.sh -t ${threshold} -f ${server_file}


build: test
	@echo "Building..."
	@go build -o ${server_output} ${server_file}

test: verify
	@echo "Testing..."
	@go test ./...

