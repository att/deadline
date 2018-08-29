package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/pkg/browser"
	"github.com/jasonlvhit/gocron"
	"github.com/att/deadline/config"
	"github.com/att/deadline/server"
	"github.com/att/deadline/common"
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
	common.Init(os.Stdout, os.Stdout)
	flag.Parse()

	if *showVersion {
		fmt.Println("Version:\t", version)
		fmt.Println("Commit:\t\t", commit)
		fmt.Println("Built by:\t", builtby)
		os.Exit(0)
	}

	cfg, err := config.LoadConfig(*configFile)

	if err != nil {
		common.Info.Println("We couldn't load the config, using defaults. Error was", err)
		cfg = &config.DefaultConfig
	}

	server.M = server.M.Init(cfg)

 	go gocron.Every(10).Seconds().Do(server.M.EvaluateAll)
	go gocron.Start() 




	dlsvr := server.NewDeadlineServer(cfg)
	err = browser.OpenURL("http://localhost:" + cfg.Server.Port + "/app")
	common.CheckError(err)



	err = dlsvr.Start()
	if err != nil {
		common.Info.Println("Server exited with error:", err)
	}
	
}
