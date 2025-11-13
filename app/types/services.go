package types

import (
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	Name  string  `gorm:"not null"`
	Price float64 `gorm:"not null"`
}
