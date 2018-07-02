#!/bin/sh

NAME="jeffo-test"
SUCCESS=true
DETAILS="{}"

BODY="{\"name\": \"${NAME}\", \"success\":${SUCCESS}, \"details\": \"${DETAILS}\"}"

curl -X POST http://localhost:8081/ -d '{"name": "jeff-test","success":true}'
