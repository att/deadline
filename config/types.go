package config

type Config struct {
	FileConfig  FileConfig  `yaml:"fileconfig"`
	DBConfig    DBConfig    `yaml:"dbconfig"`
	DAO         string      `yaml:"dao"`
	Server      ServConfig  `yaml:"serverconfig"`
	EmailConfig EmailConfig `yaml:"emailconfig"`
}

type FileConfig struct {
	Directory string `yaml:"directory"`
}
type DBConfig struct {
	ConnectionString string `yaml:"connection_string"`
}

type ServConfig struct {
	Port string `yaml:"port"`
}

type HandlerConfig struct {
	EmailConfig EmailConfig `yaml:"emailconfig"`
}

type EmailConfig struct {
	From      string `yaml:"from"`
	RelayHost string `yaml:"relay_host"`
	RelayPort int    `yaml:"relay_port"`
}
