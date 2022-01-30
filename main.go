package main

import (
	"encoding/json"
	"hiveon_monitoring/config"
	"hiveon_monitoring/logger"
	"hiveon_monitoring/psql"
	"hiveon_monitoring/scheduler"
	"io/ioutil"
)

func init() {
	if err := readConfigFile("config/local.json", config.AppConfig); err != nil {
		logger.Logging.Info("[init]: failed to read config file")
		panic(err.Error())
	}

	if err := psql.Init(&config.AppConfig.Postgres); err != nil {
		logger.Logging.Info("[init]: failed to init postgres db")
		panic(err.Error())
	}

	// if err := psql.CreateTables(); err != nil {
	// 	logger.Logging.Info("[init]: failed to create tables on db")
	// 	panic(err.Error())
	// }

	if err := logger.Init(); err != nil {
		logger.Logging.Info("[init]: failed to init logger")
		panic(err.Error())
	}
}

func main() {
	defer cleanup()
	s := scheduler.Init()
	s.Run()
}

func readConfigFile(path string, conf *config.Config) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, conf); err != nil {
		return err
	}

	return nil
}

func cleanup() {
	psql.Close()
	logger.Close()
}
