package config

import (

	"testing"

	"github.com/stretchr/testify/assert"
)

var good = Config{
	DAO: "file",
	FileConfig: FileConfig{
		Directory: "../server/",
	},
	DBConfig: DBConfig{
		ConnectionString: "N/A",
	},
	Server: ServConfig{
		Port: "8081",
	},
}

func TestGoodConfig(test *testing.T) {
	g, err := LoadConfig("testdata/goodfile.toml")
	assert.Nil(test, err, "Could not load a good file")
	assert.Equal(test, g, &good)

}

func TestBadConfig(test *testing.T) {
	g, err := LoadConfig("testdata/badfile.toml")
	assert.NotNil(test, err, "Loaded a bad file")
	assert.Equal(test, &DefaultConfig, g)
}
