package main

import (
	"fmt"

	"egbitbucket.dtvops.net/deadline/server"
)

func main() {
	dlsvr := server.NewDeadlineServer()
	err := dlsvr.Start()
	if err != nil {
		fmt.Println("Server exited with error:", err)
	}
}
