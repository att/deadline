package config

import "net/http"

type Config struct {
	FileConfig 	FileConfig		`toml:"fileconfig"`
	DBConfig	DBConfig		`toml:"dbconfig"`
	DAO			string			`toml:"dao"`
	Server 		http.Server		`toml:"server"`			//fix pls

}

type FileConfig struct {
	Port 		string			`toml:"port"`
	Path		string 			`toml:"path"`
	Host 		string 			`toml:"host"`
	
}

type DBConfig struct {
	Port 		string			`toml:"port"`
	Name		string 			`toml:"name"`
	Username	string 			`toml:"username"`
	Password	string 			`toml:"password"`
	Host 		string 			`toml:"host"`
	Path		string 			`toml:"path"`

}

