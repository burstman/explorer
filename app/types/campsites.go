package types

import (
	"time"

	"gorm.io/gorm"
)

type CampSite struct {
	ID            int            `gorm:"primaryKey"`
	Name          string         `gorm:"column:title"`
	Buses         []BuseType     `gorm:"many2many:camp_buses;"`
	Description   string         `gorm:"column:description"`
	ImageURL      string         `gorm:"column:image_url"`
	Location      string         `gorm:"column:location"`
	AvailableFrom *time.Time     `gorm:"column:available_from"`
	AvailableTo   *time.Time     `gorm:"column:available_to"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (CampSite) TableName() string {
	return "campsites"
}
