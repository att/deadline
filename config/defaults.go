

package config


var DefaultConfig = Config{
	FileConfig: DefaultFileConfig,
	Server:     DefaultServerConfig,
	DAO:        "file",
}

var DefaultFileConfig = FileConfig{
	Directory: ".",
}

var DefaultDBConfig = DBConfig{

	ConnectionString: "wearenotworkingwithdatabasesyet",
}

var DefaultEmailConfig = EmailConfig{}

var DefaultServerConfig = ServConfig{
	Port: "8081",
}