package psql

// docker run --name hiveon_monitoring -e POSTGRES_DB=postgres POSTGRES_USER=postgres -e POSTGRES_PASSWORD=trm -p 5432:5432 -d postgres

import (
	"fmt"
	"hiveon_monitoring/config"
	"hiveon_monitoring/entities"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var Session *gorm.DB

func GetTables() []interface{} {
	tables := []interface{}{
		&entities.OfflineWorker{},
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

func GetWorker(name string) (*entities.OfflineWorker, error) {
	var worker entities.OfflineWorker
	err := Session.Model(&worker).Where("name = ? AND downtime_length = nil", name).Find(&worker).Error
	return &worker, err
}

func UpdateWorker(worker *entities.OfflineWorker) error {
	err := Session.Model(&worker).Where("id = ?", worker.Id).Update(&worker).Error
	return err
}

func DeleteWorker(worker *entities.OfflineWorker) error {
	err := Session.Model(&worker).Where("id = ?", worker.Id).Delete(&worker).Error
	return err
}
