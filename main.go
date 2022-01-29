package main

import (
	"encoding/json"
	"hiveon_monitoring/config"
	"hiveon_monitoring/psql"
	"hiveon_monitoring/scheduler"
	"io/ioutil"
)

func init() {
	if err := readConfigFile("config/local.json", config.AppConfig); err != nil {
		panic(err.Error())
	}

	if err := psql.Init(&config.AppConfig.Postgres); err != nil {
		panic(err.Error())
	}

	if err := psql.CreateTables(); err != nil {
		panic(err.Error())
	}
}

func main() {
	scheduler.RunScheduler()
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
