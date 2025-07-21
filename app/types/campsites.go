package types

import (
	"time"

	"gorm.io/gorm"
)

type CampSite struct {
	ID            int            `gorm:"primaryKey;autoIncrement"`
	Name          string         `gorm:"column:title"`
	Buses         []BuseType     `gorm:"many2many:campsite_buses;joinForeignKey:CampsiteID;joinReferences:BusTypeID"`
	Description   string         `gorm:"column:description"`
	ImageURL      string         `gorm:"column:image_url"`
	Location      string         `gorm:"column:location"`
	Price         float64        `gorm:"column:price"`
	AvailableFrom *time.Time     `gorm:"column:available_from"`
	AvailableTo   *time.Time     `gorm:"column:available_to"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (CampSite) TableName() string {
	return "campsites"
}

type CampsiteBus struct {
	CampsiteID int `gorm:"primaryKey"`
	BusTypeID  int `gorm:"primaryKey"`
	Quantity   int

	Campsite CampSite `gorm:"foreignKey:CampsiteID"`
	BusType  BuseType `gorm:"foreignKey:BusTypeID"`
}
