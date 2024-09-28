package models

type DatabaseConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Dbhost   string `yaml:"dbhost"`
	Dbport   string `yaml:"dbport"`
	Dbname   string `yaml:"dbname"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"port"`
}

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}
