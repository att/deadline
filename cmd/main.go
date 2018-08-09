package main

import (
	"flag"

	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"

	"egbitbucket.dtvops.net/deadline/config"
	"egbitbucket.dtvops.net/deadline/schedule"
	"egbitbucket.dtvops.net/deadline/server"
)

var (
	configFile  *string = flag.String("config.file", "config.toml", "The Config file this binary is to run with")
	showVersion *bool   = flag.Bool("version", false, "Show the version of this binary")

	version string = NOT_DEFINED
	commit  string = NOT_DEFINED
	builtby string = NOT_DEFINED
)

const (
	NOT_DEFINED string = "Not defined"
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
		fmt.Println("We couldn't load the config, using defaults. Error was", err)
		cfg = &config.DefaultConfig
	}

	spew.Dump(cfg)

	server.Fd = schedule.NewScheduleDAO(cfg)

	server.M = schedule.NewManager()

	dlsvr := server.NewDeadlineServer(cfg)

	err = dlsvr.Start()
	if err != nil {
		fmt.Println("Server exited with error:", err)
	}
}
