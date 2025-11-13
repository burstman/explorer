package types

import "gorm.io/gorm"

type LandingVideo struct {
	gorm.Model
	Enabled bool   `gorm:"not null;default:true"`
	URL     string `gorm:"not null"`
}
