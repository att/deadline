package config

import (
//"github.com/davecgh/go-spew/spew"
"github.com/stretchr/testify/assert"
"testing" 
)

var good = Config{
	DAO: "file",
	FileConfig: FileConfig{
		Directory: "../server/",
	},
	DBConfig: DBConfig{
		ConnectionString: "somethintoo",
	},
	Server: ServConfig{
		Port: "8081",
	},


}

func TestGoodConfig(test *testing.T) {
	g, err := LoadConfig("goodfile.toml")
	assert.Nil(test, err, "Could not load a good file")
	assert.Equal(test,g,&good)
	//should print goodfile struct

}

func TestBadConfig(test *testing.T) {
	g, err := LoadConfig("badfile.toml")
	assert.NotNil(test, err, "Loaded a bad file")
	assert.Equal(test,&DefaultConfig,g) 
	//should print default
}