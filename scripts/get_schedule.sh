#!/bin/sh


server='localhost'
port=$2
schedule=$1

curl -v -X GET http://${server}:${port}/api/v1/schedule?name=${schedule} | jq '.' 
