package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/att/deadline/config"
	"github.com/att/deadline/server"
	"github.com/jasonlvhit/gocron"
)

var (
	configFile  = flag.String("config.file", "config.toml", "The Config file this binary is to run with")
	showVersion = flag.Bool("version", false, "Show the version of this binary")

	version = NotDefined
	commit  = NotDefined
	builtby = NotDefined
)

const (
	// NotDefined is a simple const string for things that are not defined
	NotDefined = "Not defined"
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println("Version:\t", version)
		fmt.Println("Commit:\t\t", commit)
		fmt.Println("Built by:\t", builtby)
		os.Exit(0)
	}

	cfg, err := config.LoadConfig(*configFile)

	if err != nil {
		cfg = &config.DefaultConfig
	}

	//server.M = server.M.Init(cfg)

	//go gocron.Every(10).Seconds().Do(server.M.EvaluateAll)
	go gocron.Start()

	dlsvr := server.NewDeadlineServer(cfg)

	err = dlsvr.Start()
	if err != nil {

	}
}
