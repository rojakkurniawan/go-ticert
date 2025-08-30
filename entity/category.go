package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	EventID   uuid.UUID      `json:"event_id" gorm:"type:char(36);not null"`
	Name      string         `json:"name" gorm:"type:varchar(255);not null"`
	Price     float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	EventDate time.Time      `json:"event_date" gorm:"type:datetime;not null"`
	Quantity  int            `json:"quantity" gorm:"type:int;not null"`
	Status    string         `json:"status" gorm:"type:enum('available','sold');not null;default:'available'"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Event *Event `json:"event" gorm:"foreignKey:EventID"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
