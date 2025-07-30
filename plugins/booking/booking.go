package booking

import (
	"explorer/app/db"
	"explorer/app/types"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/anthdm/superkit/kit"
)

func HandelCreateBooking(kit *kit.Kit) error {

	
	//  Parse form
	if err := kit.Request.ParseForm(); err != nil {
		return err
	}

	// Parse scalar fields
	auth := kit.Auth().(types.AuthUser)
	campID, err := strconv.Atoi(kit.Request.FormValue("campID"))
	if err != nil {
		return err
	}
	totalPrice, err := strconv.ParseFloat(kit.Request.FormValue("totalPrice"), 64)
	if err != nil {
		return err
	}
	paymentMethod := kit.Request.FormValue("payment_method")
	specialRequest := kit.Request.FormValue("specialRequest")

	// Booking object
	booking := types.Bookings{
		UserID:         auth.UserID,
		CampID:         uint(campID),
		SpecialRequest: specialRequest,
		TotalPrice:     totalPrice,
		Status:         "pending",
		PaymentStatus:  "unpaid",
		PaymentMethod:  paymentMethod,
	}

	// Parse guests
	guestsCount, err := strconv.Atoi(kit.Request.FormValue("guestsCount"))
	if err != nil {
		return err
	}
	for i := 0; i < guestsCount; i++ {
		guest := types.Guest{
			FirstName: kit.Request.FormValue(fmt.Sprintf("guests[%d][first_name]", i)),
			LastName:  kit.Request.FormValue(fmt.Sprintf("guests[%d][last_name]", i)),
			CIN:       kit.Request.FormValue(fmt.Sprintf("guests[%d][cin]", i)),
		}
		booking.Guests = append(booking.Guests, guest)
	}

	// Parse services
	for k, v := range kit.Request.Form {
		if strings.HasPrefix(k, "service[") {
			idStr := strings.TrimSuffix(strings.TrimPrefix(k, "service["), "]")
			serviceID, _ := strconv.Atoi(idStr)
			qty, err := strconv.Atoi(v[0])
			if err != nil {
				return err
			}
			if qty > 0 {
				booking.Services = append(booking.Services, types.BookingService{
					ServiceID: uint(serviceID),
					Quantity:  qty,
				})
			}
		}
	}

	// Create booking with nested associations
	if err := db.Get().Create(&booking).Error; err != nil {

		return err
	}

	return kit.Redirect(http.StatusSeeOther, "/")
}
