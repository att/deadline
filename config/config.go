package config

import (
	"errors"
	"os"

	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func validateConfig(c Config) error {

	if (c.DAO == "") || (c.DAO != "DB" && c.DAO != "file") {
		return errors.New("DAO not specified, used default")
	}
	return nil
}

func LoadConfig(filename string) (*Config, error) {

	var config = &Config{}

	if data, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else {
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, err
		} else {
			return config, nil
		}
	}

}

func checkMissingConfigs(c *Config) {
	switch c.DAO {
	case "DB":
		if c.DBConfig.ConnectionString == "" {
			c.DBConfig = DefaultDBConfig
		}
		break
	case "file":
		if c.FileConfig.Directory == "" {
			c.FileConfig = DefaultFileConfig
		}
		if _, err := os.Stat(c.FileConfig.Directory); os.IsNotExist(err) {
			c.FileConfig = DefaultFileConfig
		}
	}
}
