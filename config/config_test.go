package config

import (
"github.com/davecgh/go-spew/spew"
"github.com/stretchr/testify/assert"
"testing" 
)

func TestGoodConfig(test *testing.T) {
	g, err := LoadConfig("goodfile.toml")
	assert.Nil(test, err, "Could not load a good file")
	spew.Dump(g) 

}

func TestBadConfig(test *testing.T) {
	_, err := LoadConfig("badfile.toml")
	assert.NotNil(test, err, "Loaded a bad file")
}