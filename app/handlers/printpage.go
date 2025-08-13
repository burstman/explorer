package handlers

import (
	"explorer/app/db"
	"explorer/app/types"
	"explorer/app/views/landing"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/anthdm/superkit/kit"
)

func HandlePrintBookings(kit *kit.Kit) error {
	if err := kit.Request.ParseForm(); err != nil {
		return fmt.Errorf("parsing error in HandlePrintBookings: %v", err)
	}

	idsParam := kit.Request.FormValue("ids") // "1,2,3"
	log.Println("id print", idsParam)
	if idsParam == "" {
		return fmt.Errorf("parsing ids is empty in HandlePrintBookings")
	}

	// Split the comma-separated IDs into a slice
	idsStr := strings.Split(idsParam, ",") // ["1", "2", "3"]

	var ids []int
	for _, idStr := range idsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("invalid ID: %s : %v", idStr, err)
		}
		ids = append(ids, id)
	}

	// Fetch users with bookings and guests preloaded
	var users []types.User
	err := db.Get().
		Preload("Bookings").
		Preload("Bookings.Guests").
		Where("id IN ?", ids).
		Where("role <> ?", "admin").
		Find(&users).Error

	if err != nil {
		return fmt.Errorf("db query error in HandlePrintBookings: %v", err)
	}

	// Now users contains all the selected users with their bookings and guests loaded.
	// You can render or return as needed:
	return kit.Render(landing.PagePrintBookingSelectedUser(users))
}
