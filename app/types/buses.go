package types

import (
	"time"

	"gorm.io/gorm"
)

type BuseType struct {
	ID        int            `gorm:"primaryKey"`
	Name      string         `gorm:"column:name"`
	Capacity  int            `gorm:"column:capacity"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (BuseType) TableName() string {
	return "bus_types"
}
