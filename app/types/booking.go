package types

import (
	"time"
)

type Bookings struct {
	ID             int `gorm:"primaryKey"`
	UserID         int
	CampID         int
	SpecialRequest string
	TotalPrice     float64
	Status         string
	PaymentStatus  string
	PaymentMethod  string
	Guests         []Guest
	Services       map[int]int
	CreatedAt      time.Time
}

type Guest struct {
	ID        int `gorm:"primaryKey"`
	FirstName string
	LastName  string
	CIN       string
}
