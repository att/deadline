package main

import (
	"github.com/davecgh/go-spew/spew"
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"os"
	"egbitbucket.dtvops.net/deadline/config"
	"egbitbucket.dtvops.net/deadline/schedule"
	"egbitbucket.dtvops.net/deadline/server"
)
var c *config.Config
func main() {
	c = &config.DefaultConfig
	var arg string
	
	if len(os.Args) >= 2 {
		arg = os.Args[1]
		var err error
		c, err = config.LoadConfig(arg)
		if err != nil {
			fmt.Println("We couldn't load the config!")
		
		} 
	}

	spew.Dump(c)
	 //our location for the config 

	server.Fd = schedule.NewScheduleDAO(c)

	server.M = schedule.NewManager()

	dlsvr := server.NewDeadlineServer(c)

	err := dlsvr.Start()
	if err != nil {
		fmt.Println("Server exited with error:", err)
	}
}
