package config

import (
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// FileStorage is the static configuration string for file storage
	FileStorage string = "file"

	// DBStorage is the static configuration string for database storage
	DBStorage string = "db"

	// DefaultPort is the default server port
	DefaultPort string = "8080"

	// DefaultEvalString is the default evaluation time as a string
	DefaultEvalString string = "5m"

	// DefaultEvalDuration is the default evaulation time as a Duration
	DefaultEvalDuration time.Duration = time.Minute * 5
)

// Config represents the configuration struct for the entire deadline application
type Config struct {
	FileConfig  FileConfig        `yaml:"file_config"`
	DBConfig    DBConfig          `yaml:"db_config"`
	Storage     string            `yaml:"storage"`
	EvalTime    string            `yaml:"eval_timing"`
	Server      ServerConfig      `yaml:"server_config"`
	EmailConfig EmailConfig       `yaml:"email_config"`
	Logconfig   map[string]string `yaml:"logs"`
	loggers     map[string]*logrus.Logger
	modLock     sync.RWMutex
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
		Port: DefaultPort,
	},
	EvalTime:  DefaultEvalString,
	Storage:   FileStorage,
	loggers:   make(map[string]*logrus.Logger),
	Logconfig: make(map[string]string),
}
