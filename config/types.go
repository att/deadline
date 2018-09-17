package config

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

// Config represents the configuration struct for the entire deadline application
type Config struct {
	FileConfig  FileConfig        `yaml:"fileconfig"`
	DBConfig    DBConfig          `yaml:"dbconfig"`
	DAO         string            `yaml:"dao"`
	Server      ServerConfig      `yaml:"serverconfig"`
	EmailConfig EmailConfig       `yaml:"emailconfig"`
	Logconfig   map[string]string `yaml:"logs"`
	loggers     map[string]*logrus.Logger
	logLock     sync.Mutex
}

// FileConfig is the configuration type for file storage
type FileConfig struct {
	Directory string `yaml:"directory"`
}

// DBConfig is the configuration type for relational database storage
type DBConfig struct {
	ConnectionString string `yaml:"connection_string"`
}

// ServerConfig is the configuration type for the deadline server
type ServerConfig struct {
	Port string `yaml:"port"`
}

// HandlerConfig is the configuration type for handlers
type HandlerConfig struct {
	EmailConfig EmailConfig `yaml:"emailconfig"`
}

// EmailConfig is the configuration type for handlers that email
type EmailConfig struct {
	From      string `yaml:"from"`
	RelayHost string `yaml:"relay_host"`
	RelayPort int    `yaml:"relay_port"`
}

// DefaultConfig is the default configuration
var DefaultConfig = Config{
	FileConfig: FileConfig{
		Directory: os.TempDir() + "/deadline",
	},
	Server: ServerConfig{
		Port: "8080",
	},
	DAO:       "file",
	loggers:   make(map[string]*logrus.Logger),
	Logconfig: make(map[string]string),
}
