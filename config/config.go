package config

import (
	"github.com/sirupsen/logrus"

	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var levelLookup = map[string]logrus.Level{
	"warn":  logrus.WarnLevel,
	"info":  logrus.InfoLevel,
	"debug": logrus.DebugLevel,
	"error": logrus.ErrorLevel,
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

// LoadConfig loads the configuration based on the input file. Errors can occur for various
// i/o or marshalling related reasons. Defaults will be returned for primitive types, like strings.
func LoadConfig(filename string) (*Config, error) {
	var err error
	var config = &Config{}
	var data []byte

	if data, err = ioutil.ReadFile(filename); err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	if config.Logconfig == nil {
		config.Logconfig = make(map[string]string)
	}

	return config, nil

}

// GetLogger gets a logger for a particular package or sub-component. Thread safe, but currently locks
// pretty aggressively, so one should only call at the package/sub-component level, not like, per function call.
func (c *Config) GetLogger(name string) *logrus.Logger {
	var logger *logrus.Logger
	var found bool

	c.logLock.Lock() //locking strategy probably a bit aggressive
	defer c.logLock.Unlock()

	if logger, found = c.loggers[name]; !found {

		logger := logrus.New()
		logger.SetLevel(c.getLoggerLevel(name))
		c.loggers[name] = logger

		return logger
	}

	return logger
}

func (c *Config) getLoggerLevel(name string) logrus.Level {
	if lvl, found := c.Logconfig[name]; !found {
		return logrus.InfoLevel
	} else if level, found := levelLookup[lvl]; !found {
		return logrus.InfoLevel
	} else {
		return level
	}
}
