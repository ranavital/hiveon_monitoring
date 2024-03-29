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
	Postgres          PsqlConfig `json:"postgres"`
	HiveonWorkersPath string     `json:"hiveon_workers_url"`
	LoggerPath        string     `json:"logger_path"`
	TgToken           string     `json:"tg_token"`
	TgChatId          string     `json:"tg_chat_id"`
	AllWorkers        []string   `json:"all_workers"`
}
