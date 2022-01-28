package psql

import (
	"example/web-service-gin/config"
	"example/web-service-gin/entities"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var Session *gorm.DB

func GetTables() []interface{} {
	tables := []interface{}{
		&entities.Album{},
	}
	return tables
}

func Init(conf *config.PsqlConfig) error {
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", conf.Host, conf.Port, conf.User, conf.Database, conf.Password))
	if err != nil {
		return err
	}

	Session = db
	return nil
}

func CreateTables() error {
	return Session.AutoMigrate(GetTables()...).Error
}

func Close() error {
	return Session.Close()
}
