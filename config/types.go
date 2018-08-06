package config

//import "net/http"

import "egbitbucket.dtvops.net/deadline/common"

type Config struct {
	FileConfig  FileConfig  `toml:"fileconfig"`
	DBConfig    DBConfig    `toml:"dbconfig"`
	DAO         string      `toml:"dao"`
	Server      ServConfig  `toml:"serverconfig"`
	EmailConfig EmailConfig `toml:"emailconfig"`
}

type FileConfig struct {
	Directory string `toml:"directory"`
}
type DBConfig struct {
	Connection_String string `toml:"connection_string"`
}

type ServConfig struct {
	Port string `toml:"port"`
}

type EmailConfig struct {
	From       string `toml:"from"`
	Relay_Host string `toml:"relay_host"`
	Relay_Port int    `toml:"relay_port"`
}
