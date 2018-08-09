package config

//import "net/http"

//import "egbitbucket.dtvops.net/deadline/common"

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
	ConnectionString string `toml:"connection_string"`
}

type ServConfig struct {
	Port string `toml:"port"`
}

type HandlerConfig struct {
	EmailConfig EmailConfig `toml:"emailconfig"`
	//hipchat, slack, etc.
	
}

type EmailConfig struct {
	From       string `toml:"from"`
	RelayHost string `toml:"relay_host"`
	RelayPort int    `toml:"relay_port"`
}
