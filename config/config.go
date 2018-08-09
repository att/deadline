package config

import (
	"errors"
	"os"
	"github.com/BurntSushi/toml"
	"egbitbucket.dtvops.net/deadline/common"
)

func validateConfig(c Config) error {

	if (c.DAO == "") || (c.DAO != "DB" && c.DAO != "file") {
		return errors.New("DAO not specified, used default")
	}
	return nil
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
	if  err != nil {
		return &DefaultConfig,err 
	}
	checkMissingConfigs(&conf)
	return &conf, nil
}

func checkMissingConfigs(c *Config) {
	switch c.DAO {
	case "DB":
		if c.DBConfig.ConnectionString == "" {
			common.Info.Println("No DB specified.")
			c.DBConfig = DefaultDBConfig
		}
		break
	case "file":
		if c.FileConfig.Directory == "" {
			common.Info.Println("No directory specified.")
			c.FileConfig = DefaultFileConfig
		}
		if _, err := os.Stat(c.FileConfig.Directory); os.IsNotExist(err) {
			common.Info.Println("Given directory doesn't exist.")
			c.FileConfig = DefaultFileConfig
		}
	}
}

