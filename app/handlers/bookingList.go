package handlers

import (
	"explorer/app/db"
	"explorer/app/types"
	"explorer/app/views/landing"
	"explorer/plugins/booking"
	"log"
	"net/http"
	"strconv"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func HandelBooklist(kit *kit.Kit) error {

	var users []types.User

	err := db.Get().
		Where("role != ?", "admin").
		Preload("Bookings.Guests").
		Preload("Bookings.Services.Service").
		Preload("Bookings.Camp").
		Find(&users).Error
	if err != nil {
		log.Println("failed to fetch users and bookings:", err)
		return err
	}

	return RenderWithLayout(kit, landing.BookingListAdmin(users))
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

	var users []types.User
	err = db.Get().
		Where("role != ?", "admin").
		Preload("Bookings", func(db *gorm.DB) *gorm.DB {
			return db.Order("bookings.created_at DESC").
				Preload("Guests").
				Preload("Services.Service").
				Preload("Camp")
		}).
		Find(&users).Error
	if err != nil {
		log.Println("failed to fetch users with bookings:", err)
		return err
	}

	return kit.Render(landing.BookingTableRows(users))
}

func EditBooking(kit *kit.Kit) error {
	strID := chi.URLParam(kit.Request, "id")
	bookingID, err := strconv.Atoi(strID)
	if err != nil {
		log.Println(err)
		return err
	}

	var booking types.Bookings
	err = db.Get().
		Preload("Guests").
		Preload("Services.Service").
		Preload("Camp").
		Preload("User").
		First(&booking, "id = ?", bookingID).Error
	if err != nil {
		return err
	}
	var camps []types.CampSite
	err = db.Get().
		Order("created_at asc").
		Find(&camps).Error
	if err != nil {
		return err
	}

	// Get available services
	var availableServices []types.Service
	err = db.Get().
		Order("created_at asc").
		Find(&availableServices).Error
	if err != nil {
		log.Println("Failed to fetch services:", err)
		return err
	}
	return kit.Render(landing.EditBookingModal(booking, camps, availableServices))

}

func BookingShowDetail(kit *kit.Kit) error {
	strID := chi.URLParam(kit.Request, "id")

	bookingID, err := strconv.Atoi(strID)
	if err != nil {
		return err
	}

	var booking types.Bookings
	err = db.Get().
		Where("id = ?", bookingID).
		Preload("Guests").
		Preload("Services.Service").
		Preload("User").
		Preload("Camp").
		First(&booking).Error
	if err != nil {
		log.Println("failed to fetch booking:", err)
		return err
	}

	//log.Println("booking detail:", booking)

	return kit.Render(landing.BookingDetailModal(booking))
}

func BookingAdmin(kit *kit.Kit) error {
	strUserID := chi.URLParam(kit.Request, "user_id")

	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		return err
	}

	var user types.User

	err = db.Get().
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		return err
	}

	services, err := GetAllAvailableServices()
	if err != nil {
		return err
	}

	camps, err := GetAllAvailableCamps()
	if err != nil {
		return err
	}

	return kit.Render(landing.BookingAdminCreateModal(user, services, camps))
}

func GetAllAvailableServices() ([]types.Service, error) {
	var services []types.Service
	err := db.Get().Find(&services).Error
	if err != nil {
		return nil, err
	}
	return services, nil
}

func GetAllAvailableCamps() ([]types.CampSite, error) {
	var camps []types.CampSite
	err := db.Get().Find(&camps).Error
	if err != nil {
		return nil, err
	}
	return camps, nil
}

func AdminBookingAdd(kit *kit.Kit) error {

	var addBooking types.Bookings
	//  Parse form
	if err := kit.Request.ParseForm(); err != nil {
		return err
	}

	strUserID := chi.URLParam(kit.Request, "userID")
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		return err
	}

	strCampID := kit.Request.FormValue("camp_id")

	campID, err := strconv.Atoi(strCampID)
	if err != nil {
		return err
	}

	specialRequest := kit.Request.FormValue("specialRequest")
	strTotalPrice := kit.Request.FormValue("totalPrice")

	totalPrice, err := strconv.ParseFloat(strTotalPrice, 64)
	if err != nil {
		return err
	}

	guests, err := booking.GuestParsing(kit)
	if err != nil {
		return err
	}
	bookingServices, err := booking.BookingServiceParsing(kit)
	if err != nil {
		return err
	}

	addBooking = types.Bookings{
		UserID:         uint(userID),
		CampID:         uint(campID),
		SpecialRequest: specialRequest,
		TotalPrice:     totalPrice,
		Guests:         guests,
		Services:       bookingServices,
		Status:         types.StatusBooked,
		PaymentStatus:  types.StatusCompleted,
		PaymentMethod:  types.CashPaymentMethod,
	}

	if err := db.Get().Create(&addBooking).Error; err != nil {
		return err
	}

	return kit.Redirect(http.StatusSeeOther, "/admin/booking/list")
}
