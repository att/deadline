package config

import (
	"errors"
	"github.com/BurntSushi/toml"
	"os"
)

var DefaultConfig = Config{
	FileConfig: DefaultFileConfig,
}

var DefaultFileConfig = FileConfig{
	Port: "8080",
	Path: "goodfile.toml",
	Host: "localhost:",
}

var DefaultDBConfig = DBConfig{

	Username: "user",
	Password: "pw",
}

var DefaultEmailConfig = EmailConfig{}

var DefaultServerConfig = ServerConfig{
	Port: 8081,
}

func validateConfig(c Config) error {
	if (c.FileConfig.Port == "" && c.FileConfig.Path == "") && (c.DBConfig.Port == "" && c.DBConfig.Path == "") {
		return errors.New("no valid configs, using a default config")
	}
	if c.DAO == "" {
		return errors.New("DAO not specified")
	}
	return nil
	//db checks later
}

func LoadConfig(file string) (*Config, error) {

	_, err := os.Stat(file)
	if err != nil {
		return nil, err
	}

	var conf Config
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		return &DefaultConfig, err
	}
	err = validateConfig(conf)
	if err != nil {

		return &DefaultConfig, errors.New("the struct was empty")
	}
	return &conf, nil
}
