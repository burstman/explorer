package types

import "gorm.io/gorm"

type CarouselImage struct {
	gorm.Model        // ID, CreatedAt, UpdatedAt, DeletedAt
	URL        string `gorm:"not null"`
	Position   int    `gorm:"default:0"`
}
