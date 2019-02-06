package config

import (
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

var good = Config{
	Storage: FileStorage,
	FileConfig: FileConfig{
		Directory: "../server/",
	},
	DBConfig: DBConfig{
		ConnectionString: "N/A",
	},
	Server: ServerConfig{
		Port: "8081",
	},
	Logconfig: make(map[string]string),
	loggers:   make(map[string]*logrus.Logger),
}

func TestGoodConfig(test *testing.T) {
	cfg, err := LoadConfig("testdata/goodfile.yml")
	assert.Nil(test, err, "Could not load a good file")
	assert.Equal(test, &good, cfg)
	assert.NotNil(test, cfg.Logconfig)
}

func TestCantFindConfig(test *testing.T) {
	cfg, err := LoadConfig("testdata/cantfindfile.yml")
	assert.NotNil(test, err, "Loaded a bad file")
	assert.Nil(test, cfg)
}

func TestGetLogger(test *testing.T) {
	cfg, err := LoadConfig("testdata/file_with_loggers.yml")
	assert.Nil(test, err, "Could not load a good file")
	assert.NotNil(test, cfg.loggers)

	managerLogger := cfg.GetLogger("manager")
	assert.NotNil(test, managerLogger)
	assert.True(test, managerLogger.Level == logrus.DebugLevel)
}
