package status

import (
	"errors"
	"explorer/app/db"
	"explorer/app/handlers"
	"explorer/app/types"

	"github.com/anthdm/superkit/kit"
	"gorm.io/gorm"
)

func BookingHandler(kit *kit.Kit) error {
	user := kit.Auth().(types.AuthUser)
	userID := user.GetUserID()

	var booking types.Bookings
	err := db.Get().
		Where("user_id = ?", userID).
		Preload("Guests").
		Preload("Services.Service").
		Order("created_at DESC").
		First(&booking).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return handlers.RenderWithLayout(kit, BookingStatusPage(types.Bookings{}, types.CampSite{}))
	}

	var camp types.CampSite
	if err := db.Get().First(&camp, booking.CampID).Error; err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return handlers.RenderWithLayout(kit, BookingStatusPage(booking, camp))
}
