package types

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	ID        int            `gorm:"primaryKey"`
	Name      string         `gorm:"not null"`
	Price     float64        `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Service) TableName() string {
	return "service"
}
