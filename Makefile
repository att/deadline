
DESTDIR ?= ~

server_file	        ?=	"cmd/deadline-server.go"
server_output		?=	"deadline-server"

threshold		?=	"10"
client_pkg		?=	"client/"

all: verify test build

verify: 
	@go vet ./...
	@sh check_complexity.sh -t ${threshold} -f ${appd_utils_file}
	@sh check_complexity.sh -t ${threshold} -f ${confluence_utils_file}
	@sh check_complexity.sh -t ${threshold} -f ${splunk_utils_file}
	@golint ./... | grep -v "testfiles\|templates\|ALL_CAPS\|vendor"


build: test
	@echo "Building..."
	@go build -o ${appd_utils_output} ${appd_utils_file}
	@go build -o ${confluence_utils_output} ${confluence_utils_file}
	@go build -o ${splunk_utils_output} ${splunk_utils_file}

test: verify
	@echo "Testing..."
	@go test ./...

install: all
	install -m 0755 -d $(DESTDIR)/usr/bin
	install -m 0755 ${appd_utils_output} $(DESTDIR)/usr/bin
	install -m 0755 ${confluence_utils_output} $(DESTDIR)/usr/bin
	install -m 0755 ${splunk_utils_output} $(DESTDIR)/usr/bin
