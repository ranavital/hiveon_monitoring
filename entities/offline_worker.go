package entities

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type OfflineWorker struct {
	Id             uuid.UUID  `json:"id" gorm:"primary_key;column:id"`
	CreatedAt      *time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt      *time.Time `json:"updated_at" gorm:"not null"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	Name           string     `json:"name" gorm:"not null"`
	LastAlertTime  *time.Time `json:"last_alert_time"`
	DowntimeLength *time.Duration `json:"downtime_length"`
}
