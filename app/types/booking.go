package types

import "time"

type Booking struct {
	ID             int `gorm:"primaryKey"`
	UserID         int
	GuestCount     int
	Nights         int
	SpecialRequest string
	Services       []Service `gorm:"type:jsonb" json:"services"`
	CreatedAt      time.Time
}
