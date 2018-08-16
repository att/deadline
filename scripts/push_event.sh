#!/bin/sh

NAME="jeffo-test"
SUCCESS=true
DETAILS="{}"

BODY="{\"name\": \"${NAME}\", \"success\":${SUCCESS}, \"details\": \"${DETAILS}\"}"

file=$1
HOST=localhost
PORT=$2

curl -v -X PUT http://${HOST}:${PORT}/api/v1/event -d @../server/testdata/${file}

