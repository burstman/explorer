package handlers

import (
	"explorer/app/db"
	"explorer/app/types"
	"explorer/app/views/landing"
	"fmt"
	"strings"

	"github.com/anthdm/superkit/kit"
	"gorm.io/gorm"
)

func HandleBookingSearch(kit *kit.Kit) error {
	if err := kit.Request.ParseForm(); err != nil {
		return err
	}

	q := strings.TrimSpace(kit.Request.FormValue("q"))
	paymentStatus := kit.Request.FormValue("payment_status")
	userStatus := kit.Request.FormValue("user_status")
	payMethod := kit.Request.FormValue("payment_method")

	var users []types.User
	query := db.Get().
		Preload("Bookings.Camp").
		Preload("Bookings.Guests").
		Preload("Bookings.Services.Service").
		Where("role <> ?", "admin") // exclude admin users

	if q != "" {
		query = query.Joins("LEFT JOIN bookings ON bookings.user_id = users.id AND bookings.deleted_at IS NULL").
			Joins("LEFT JOIN campsites ON campsites.id = bookings.camp_id").
			Where(`first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR campsites.title ILIKE ?`,
				"%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%")
	}

	if paymentStatus != "" && paymentStatus != "All Payment Statuses" {
		query = query.Joins("JOIN bookings ON bookings.user_id = users.id AND bookings.deleted_at IS NULL").
			Where("bookings.payment_status = ?", strings.ToLower(paymentStatus))
	}

	if userStatus != "" && userStatus != "All User Statuses" {
		query = query.Joins("JOIN bookings b2 ON b2.user_id = users.id AND b2.deleted_at IS NULL").
			Where("b2.status = ?", strings.ToLower(userStatus))
	}

	if payMethod != "" && payMethod != "All Payment Methodes" {
		query = query.Joins("JOIN bookings b3 ON b3.user_id = users.id AND b3.deleted_at IS NULL").
			Where("b3.payment_method = ?", strings.ToLower(payMethod))
	}

	if err := query.Find(&users).Error; err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("err in HandleBookingSearch db: %v", err)
	}

	return RenderWithLayout(kit, landing.BookingTableRows(users))
}
