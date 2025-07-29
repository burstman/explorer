package types

import "gorm.io/gorm"

type Bookings struct {
	gorm.Model
	UserID         uint
	CampID         uint
	SpecialRequest string
	TotalPrice     float64
	Status         string
	PaymentStatus  string
	PaymentMethod  string

	Guests   []Guest          `gorm:"foreignKey:BookingID"`
	Services []BookingService `gorm:"foreignKey:BookingID"`
}

type Guest struct {
	gorm.Model
	BookingID uint
	FirstName string
	LastName  string
	CIN       string
}

type BookingService struct {
	gorm.Model
	BookingID uint
	ServiceID uint
	Quantity  int

	Service Service `gorm:"foreignKey:ServiceID"`
}
