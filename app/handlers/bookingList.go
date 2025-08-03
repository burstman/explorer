package handlers

import (
	"explorer/app/db"
	"explorer/app/types"
	"explorer/app/views/landing"
	"log"
	"net/http"
	"strconv"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
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

func HandelDeleteBookingList(kit *kit.Kit) error {

	strBookingID := chi.URLParam(kit.Request, "bookID")

	BookingID, err := strconv.Atoi(strBookingID)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := db.Get().Delete(&types.Bookings{}, BookingID).Error; err != nil {
		log.Println(err)
		return err
	}

	var bookings []types.Bookings
	err = db.Get().
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

	return kit.Render(landing.BookingTableRows(bookings))
}

func EditBooking(kit *kit.Kit) error {
	strID := chi.URLParam(kit.Request, "id")
	log.Println("book edit Id:", strID)

	return kit.Text(http.StatusAccepted, strID)
}
