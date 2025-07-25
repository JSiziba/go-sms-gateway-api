package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type SMSMessage struct {
	Id             string    `json:"id" gorm:"primaryKey"`
	PhoneNumber    string    `json:"phoneNumber" gorm:"type:varchar(20);not null"`
	Message        string    `json:"message" gorm:"type:varchar(500);not null"`
	CreatedAt      time.Time `json:"createdAt" gorm:"autoCreateTime"`
	DeliveryStatus string    `json:"deliveryStatus" gorm:"type:varchar(20);not null"`
}

func (entity *SMSMessage) BeforeCreate(tx *gorm.DB) (err error) {
	if entity.Id == "" {
		entity.Id = uuid.New().String()
	}

	return
}

type SMSMessageRequestDto struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Message     string `json:"message" binding:"required"`
}
