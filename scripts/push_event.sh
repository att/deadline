#!/bin/sh

NAME="jeffo-test"
SUCCESS=true
DETAILS="{}"

BODY="{\"name\": \"${NAME}\", \"success\":${SUCCESS}, \"details\": \"${DETAILS}\"}"

SERVER=localhost
PORT=8081

curl -v -X POST http://${HOST}:${PORT}/api/v1/event -d ${BODY}
