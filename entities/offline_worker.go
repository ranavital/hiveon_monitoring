package entities

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type OfflineWorker struct {
	Id             uuid.UUID  `json:"id" gorm:"primary_key;column:id"`
	CreatedAt      *time.Time `json:"created_at" gorm:"not null;column:created_at"`
	UpdatedAt      *time.Time `json:"updated_at" gorm:"not null;column:updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
	Name           string     `json:"name" gorm:"column:name"`
	LastAlertTime  *time.Time `json:"last_alert_time" gorm:"column:last_alert_time"`
	DowntimeLength string     `json:"downtime_length" gorm:"column:downtime_length"`
}

func (ow *OfflineWorker) BeforeCreate() {
	if ow.Id == uuid.Nil {
		ow.Id = uuid.NewV4()
	}
}
