package handlers

import (
	"explorer/app/types"
	"explorer/app/views/landing"

	"github.com/anthdm/superkit/kit"
)

func HandelBooklist(kit *kit.Kit) error {

	var bookinglist []types.BookingDetails
	return RenderWithLayout(kit, landing.BookingListAdmin(bookinglist))
}
