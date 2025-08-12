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
	query := db.Get().Model(&types.User{}).
		Preload("Bookings.Camp").
		Preload("Bookings.Guests").
		Preload("Bookings.Services.Service").
		Where("role <> ?", "admin").
		Joins("LEFT JOIN bookings b ON b.user_id = users.id AND b.deleted_at IS NULL").
		Joins("LEFT JOIN campsites ON campsites.id = b.camp_id")

	if q != "" {
		query = query.Where(`first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR campsites.title ILIKE ?`,
			"%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%")
	}

	if paymentStatus != "" && paymentStatus != "All Payment Statuses" {
		query = query.Where("b.payment_status = ?", strings.ToLower(paymentStatus))
	}

	if userStatus != "" && userStatus != "All User Booking Statuses" {
		query = query.Where("b.status = ?", strings.ToLower(userStatus))
	}

	if payMethod != "" && payMethod != "All Payment Methods" {
		query = query.Where("b.payment_method = ?", strings.ToLower(payMethod))
	}

	if err := query.Find(&users).Error; err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("err in HandleBookingSearch db: %v", err)
	}

	return RenderWithLayout(kit, landing.BookingTableRows(users))
}
