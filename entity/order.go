package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID         uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	CategoryID uuid.UUID      `json:"category_id" gorm:"type:char(36);not null"`
	UserID     uuid.UUID      `json:"user_id" gorm:"type:char(36);not null"`
	InvoiceID  string         `json:"invoice_id" gorm:"type:varchar(255);not null;unique"`
	Status     string         `json:"status" gorm:"type:enum('pending','paid','cancelled');not null;default:'pending'"`
	Quantity   int            `json:"quantity" gorm:"type:int;not null"`
	TotalPrice float64        `json:"total_price" gorm:"type:decimal(10,2);not null"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Category     *Category      `json:"category" gorm:"foreignKey:CategoryID"`
	User         *User          `json:"user" gorm:"foreignKey:UserID"`
	OrderDetails []*OrderDetail `json:"order_details" gorm:"foreignKey:OrderID;references:ID"`
}

type OrderDetail struct {
	ID             uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	OrderID        uuid.UUID      `json:"order_id" gorm:"type:char(36);not null"`
	TicketCode     string         `json:"ticket_code" gorm:"type:varchar(255);not null;unique"`
	FullName       string         `json:"full_name" gorm:"type:varchar(255);not null"`
	IdentityNumber string         `json:"identity_number" gorm:"type:varchar(255);not null"`
	Redeemed       bool           `json:"redeemed" gorm:"type:boolean;not null;default:false"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Order *Order `json:"order" gorm:"foreignKey:OrderID;references:ID"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

func (o *OrderDetail) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}
