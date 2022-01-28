package config

var AppConfig = &Config{}

type PsqlConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type Config struct {
	Postgres PsqlConfig `json:"postgres"`
}
