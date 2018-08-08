package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

var DefaultConfig = Config{
	FileConfig: DefaultFileConfig,
	Server:     DefaultServerConfig,
	DAO:        "file",
}

var DefaultFileConfig = FileConfig{
	Directory: ".",
}

var DefaultDBConfig = DBConfig{

	ConnectionString: "wearenotworkingwithdatabasesyet",
}

var DefaultEmailConfig = EmailConfig{}

var DefaultServerConfig = ServConfig{
	Port: "8081",
}

func validateConfig(c Config) error {

	if (c.DAO == "") || (c.DAO != "DB" && c.DAO != "file") {

		return errors.New("DAO not specified, used default")
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
		return &DefaultConfig, err
	}

	if (conf.DAO == "DB") && (conf.DBConfig.ConnectionString == "") {
		fmt.Println("No DB specified.")
		conf.DBConfig = DefaultDBConfig
	}

	if conf.DAO == "file" {
		if conf.FileConfig.Directory == "" {
			fmt.Println("No directory specified.")
			conf.FileConfig = DefaultFileConfig
		}
		if _, err := os.Stat(conf.FileConfig.Directory); os.IsNotExist(err) {
			fmt.Println("Given directory doesn't exist.")
			conf.FileConfig = DefaultFileConfig
		}
	}

	return &conf, nil
}
