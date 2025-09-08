package handlers

import (
	"explorer/app/db"
	"explorer/app/types"
	"explorer/app/views/landing"
	"explorer/plugins/booking"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	statusOptions, err := GetStatusOptions()

	if err != nil {
		return err
	}

	return RenderWithLayout(kit, landing.BookingListAdmin(users, statusOptions))
}

func GetStatusOptions() (types.StatusOptions, error) {
	var paymentStatuses []string
	var userStatuses []string

	// Distinct payment_status values
	err := db.Get().Model(&types.Bookings{}).
		Distinct("payment_status").
		Pluck("payment_status", &paymentStatuses).Error
	if err != nil {
		return types.StatusOptions{}, fmt.Errorf("err in GetStatusOptions: %v", err)
	}

	// Distinct user status values
	err = db.Get().Model(&types.Bookings{}).
		Distinct("status").
		Pluck("status", &userStatuses).Error
	if err != nil {
		return types.StatusOptions{}, fmt.Errorf("err in GetStatusOptions: %v", err)
	}

	return types.StatusOptions{
		PaymentStatuses: paymentStatuses,
		UserStatuses:    userStatuses,
	}, nil
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

func EditPostBooking(kit *kit.Kit) error {
	// Get booking ID from URL
	strBookID := chi.URLParam(kit.Request, "Bookid")
	bookingID, err := strconv.Atoi(strBookID)
	if err != nil {
		return fmt.Errorf("invalid booking ID: %w", err)
	}

	// Load existing booking
	var booking types.Bookings
	if err := db.Get().
		Preload("Guests").
		Preload("Services").
		First(&booking, bookingID).Error; err != nil {
		return fmt.Errorf("booking not found: %w", err)
	}

	// Parse form values
	if err := kit.Request.ParseForm(); err != nil {
		return fmt.Errorf("parse form: %w", err)
	}

	// Camp selection
	campID, _ := strconv.Atoi(kit.Request.FormValue("camp_id"))
	booking.CampID = uint(campID)

	// Special requests
	booking.SpecialRequest = kit.Request.FormValue("specialRequest")

	// Guests
	guestsCount, _ := strconv.Atoi(kit.Request.FormValue("guestsCount"))

	newGuests := make([]types.Guest, 0, guestsCount)
	for i := 0; i < guestsCount; i++ {
		first := kit.Request.FormValue(fmt.Sprintf("guests[%d][first_name]", i))
		last := kit.Request.FormValue(fmt.Sprintf("guests[%d][last_name]", i))
		cin := kit.Request.FormValue(fmt.Sprintf("guests[%d][cin]", i))
		newGuests = append(newGuests, types.Guest{
			FirstName: first,
			LastName:  last,
			CIN:       cin,
			BookingID: uint(bookingID),
		})
	}

	// Delete old guests and insert new ones
	if err := db.Get().Where("booking_id = ?", bookingID).Delete(&types.Guest{}).Error; err != nil {
		return fmt.Errorf("delete old guests: %w", err)
	}
	if len(newGuests) > 0 {
		if err := db.Get().Create(&newGuests).Error; err != nil {
			return fmt.Errorf("insert new guests: %w", err)
		}
	}

	// Services
	services := []types.BookingService{}
	for key, vals := range kit.Request.Form {
		if strings.HasPrefix(key, "services[") {
			// Extract service ID
			serviceIDStr := strings.TrimSuffix(strings.TrimPrefix(key, "services["), "]")
			serviceID, _ := strconv.Atoi(serviceIDStr)

			quantity, _ := strconv.Atoi(vals[0])
			if quantity > 0 {
				services = append(services, types.BookingService{
					ServiceID: uint(serviceID),
					Quantity:  quantity,
					BookingID: uint(bookingID),
				})
			}
		}
	}
	// Delete old services and insert new ones
	if err := db.Get().Where("booking_id = ?", bookingID).Delete(&types.BookingService{}).Error; err != nil {
		return fmt.Errorf("delete old services: %w", err)
	}
	if len(services) > 0 {
		if err := db.Get().Create(&services).Error; err != nil {
			return fmt.Errorf("insert new services: %w", err)
		}
	}

	// Total price
	totalPrice, _ := strconv.ParseFloat(kit.Request.FormValue("totalPrice"), 64)
	booking.TotalPrice = totalPrice

	// Status updates
	booking.Status = kit.Request.FormValue("user_status")
	booking.PaymentStatus = kit.Request.FormValue("payment_status")

	// Save booking main fields
	if err := db.Get().Model(&booking).Updates(map[string]interface{}{
		"camp_id":         booking.CampID,
		"special_request": booking.SpecialRequest,
		"total_price":     booking.TotalPrice,
		"status":          booking.Status,
		"payment_status":  booking.PaymentStatus,
	}).Error; err != nil {
		return fmt.Errorf("update booking: %w", err)
	}

	return kit.Redirect(http.StatusSeeOther, "/admin/booking/list")
}
