package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Event struct {
	ID          uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	Organizer   string         `json:"organizer" gorm:"type:varchar(255);not null"`
	Title       string         `json:"title" gorm:"type:varchar(255);not null"`
	Description string         `json:"description" gorm:"type:text;not null"`
	StartDate   time.Time      `json:"start_date" gorm:"type:datetime;not null"`
	EndDate     time.Time      `json:"end_date" gorm:"type:datetime;not null"`
	StartTime   time.Time      `json:"start_time" gorm:"type:datetime;not null"`
	EndTime     time.Time      `json:"end_time" gorm:"type:datetime;not null"`
	Location    string         `json:"location" gorm:"type:varchar(255);not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Categories []Category `json:"categories" gorm:"foreignKey:EventID"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
