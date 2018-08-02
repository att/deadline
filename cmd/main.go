package main

import (
	"fmt"
	"os"
	"egbitbucket.dtvops.net/deadline/config"
	"egbitbucket.dtvops.net/deadline/schedule"
	"egbitbucket.dtvops.net/deadline/server"
)

func main() {
	arg := os.Args[1]
	c, err := config.LoadConfig(arg)
	 //our location for the config 
	server.Fd = schedule.NewScheduleDAO(c)
	server.M = schedule.NewManager()
	dlsvr := server.NewDeadlineServer(c)
	err = dlsvr.Start()
	if err != nil {
		fmt.Println("Server exited with error:", err)
	}
}
