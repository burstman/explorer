package handlers

import (
	"explorer/app/db"
	"explorer/app/types"
	"explorer/app/views/landing"
	"log"

	"github.com/anthdm/superkit/kit"
)

func HandelBooklist(kit *kit.Kit) error {
	var bookings []types.Bookings

	// Load all bookings (excluding admins), with guests, services, and user
	err := db.Get().
		Preload("Guests").
		Preload("Services.Service").
		Preload("User").
		Preload("Camp"). // optional if you want camp name
		Joins("JOIN users ON users.id = bookings.user_id").
		Where("users.role != ?", "admin").
		Order("bookings.created_at DESC").
		Find(&bookings).Error

	if err != nil {
		log.Println("failed to fetch bookings:", err)
		return err
	}

	// Map data to BookingDetails for the frontend
	var bookinglist []types.BookingDetails
	for _, booking := range bookings {
		bookinglist = append(bookinglist, types.BookingDetails{
			BookingID: int(booking.ID),
			CampID:    int(booking.CampID),
			CampName:  booking.Camp.Title, // if Camp is preloaded
			User: types.User{
				ID:        booking.UserID,
				FirstName: booking.User.FirstName,
				LastName:  booking.User.LastName,
				Email:     booking.User.Email,
			},
			Guests: booking.Guests,
			// You can also add booking.Status, booking.PaymentStatus, etc. to BookingDetails if needed
		})
	}

	return RenderWithLayout(kit, landing.BookingListAdmin(bookinglist))
}
