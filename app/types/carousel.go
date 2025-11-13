package types

import (
	"time"

	"gorm.io/gorm"
)

type CarouselImage struct {
	gorm.Model        // ID, CreatedAt, UpdatedAt, DeletedAt
	URL        string `gorm:"not null"`
}

type CarouselConfig struct {
	ID        uint      `gorm:"primaryKey"`
	Enabled   bool      `gorm:"not null;default:true"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the exact table name in the DB
func (CarouselConfig) TableName() string {
	return "carousel_config"
}
