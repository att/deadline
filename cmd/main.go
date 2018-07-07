package main

import (
	"egbitbucket.dtvops.net/deadline/server"
)

func main() {
	dlsvr := server.NewDeadlineServer()
	dlsvr.Start()
}
