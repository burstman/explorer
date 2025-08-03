package types

import "gorm.io/gorm"

const (
	StatusPending   = "pending"
	StatusBooked    = "booked"
	StatusCancelled = "cancelled"
	StatusCompleted = "completed"
)

type Bookings struct {
	gorm.Model

	UserID         uint // ✅ needed for foreign key to User
	CampID         uint // ✅ needed for foreign key to CampSite
	SpecialRequest string
	TotalPrice     float64
	Status         string
	PaymentStatus  string
	PaymentMethod  string

	Guests   []Guest          `gorm:"foreignKey:BookingID"`
	Services []BookingService `gorm:"foreignKey:BookingID"`

	User User     `gorm:"foreignKey:UserID"` // ✅ GORM uses UserID to preload
	Camp CampSite `gorm:"foreignKey:CampID"` // ✅ GORM uses CampID to preload
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
