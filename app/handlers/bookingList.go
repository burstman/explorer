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

	// Load bookings with related data, exclude admins
	err := db.Get().
		Joins("JOIN users ON users.id = bookings.user_id").
		Where("users.role != ?", "admin").
		Preload("Guests").
		Preload("Services.Service").
		Preload("User").
		Preload("Camp").
		Order("bookings.created_at DESC").
		Find(&bookings).Error
	if err != nil {
		log.Println("failed to fetch bookings:", err)
		return err
	}

	return RenderWithLayout(kit, landing.BookingListAdmin(bookings))
}
