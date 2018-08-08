#!/bin/sh

file=$1
server='localhost'
port='6789'

curl -v -X PUT http://${server}:${port}/api/v1/schedule -d @${file}
