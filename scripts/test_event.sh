#!/bin/sh

NAME="jeffo-test"
SUCCESS=true
DETAILS="{}"

BODY="{\"name\": \"${NAME}\", \"success\":${SUCCESS}, \"details\": \"${DETAILS}\"}"

curl -v -X POST http://localhost:8081/api/v1/event -d '{"name": "jeff-test","success":true}'
