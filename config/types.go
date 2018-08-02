package config

type Config struct {
	FileConfig 	FileConfig		`toml:"fileconfig"`
	DBConfig	DBConfig		`toml:"dbconfig"`
	Port int					`toml:"port"`
	Path	string 				`toml:"path"`

}

type FileConfig struct {
	
	

}

type DBConfig struct {
	Name		string 			`toml:"name"`
	Username	string 			`toml:"username"`
	Password	string 			`toml:"password"`
	Host 		string 			`toml:"host"`

}

//defaultconfig?





/*
type WebhookConfig struct {




}

type EmailConfig struct {
	From 	string	
	To		string
	Msg 	string


}
*/