#!/bin/sh


server='localhost'
port=$2
schedule=$1

curl -X GET http://${server}:${port}/status?name=${schedule} > b.out 

