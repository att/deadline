package config
import (
"github.com/BurntSushi/toml"
"os"
"errors"
)

func validateConfig(c Config) error{
    if c.Port == 0 {
        return errors.New("Not a valid port")
    }

    if c.Path == "" {
        return errors.New("No path was given")
    }
    return nil
    //db checks later 
}

func LoadConfig(file string) (*Config,error) {

    _, err := os.Stat(file)
    if err != nil {
        return nil, err
    }

    var conf Config
    if _, err := toml.DecodeFile(file, &conf); err != nil {
        return nil, err
    }
    err = validateConfig(conf)
    if err != nil {
        return nil, errors.New("the struct was empty")
    }
	return  &conf,nil
}
