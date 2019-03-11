package config

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var (
	levelLookup = map[string]logrus.Level{
		"warn":  logrus.WarnLevel,
		"info":  logrus.InfoLevel,
		"debug": logrus.DebugLevel,
		"error": logrus.ErrorLevel,
	}

	formatter = &logrus.TextFormatter{
		FullTimestamp:    true,
		DisableTimestamp: false,
		TimestampFormat:  time.RFC3339,
	}

	globalConfig *Config

	cfgLock = sync.RWMutex{}
)

// GetConfig will return a global instance of the configuration if it has ever been loaded through
// LoadConfig. It can also return the default config if LoadConfig has never been called.
func GetConfig() *Config {
	cfgLock.RLock()
	defer cfgLock.RUnlock()

	if globalConfig == nil {
		return &DefaultConfig
	}

	return globalConfig
}

// GetEmailConfig will the EmailConfig portion of the global config object.
func GetEmailConfig() *EmailConfig {
	cfgLock.RLock()
	defer cfgLock.RUnlock()

	if globalConfig == nil {
		return &DefaultConfig.EmailConfig
	}

	return &globalConfig.EmailConfig
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

	if config.loggers == nil {
		config.loggers = make(map[string]*logrus.Logger)
	}

	cfgLock.Lock()
	defer cfgLock.Unlock()

	globalConfig = config
	return config, nil

}

// GetEvalTime is a simple facade for getting the configuration's EvalTime while parsing
// and checking for errors.
func (c *Config) GetEvalTime() time.Duration {
	if c.EvalTime == "" {
		return DefaultEvalDuration
	} else if duration, err := time.ParseDuration(c.EvalTime); err != nil {
		return duration
	}
	return DefaultEvalDuration
}

// GetLogger gets a logger for a particular package or sub-component. Thread safe, but currently locks
// pretty aggressively, so one should only call at the package/sub-component level, not like, per function call.
func (c *Config) GetLogger(name string) *logrus.Logger {
	var logger *logrus.Logger
	var found bool

	if logger, found = c.loggers[name]; !found {
		c.modLock.Lock() //locking strategy probably a bit aggressive
		defer c.modLock.Unlock()

		logger := logrus.New()
		logger.Formatter = formatter
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
