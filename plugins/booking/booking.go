package booking

import (
	"explorer/app/db"
	"explorer/app/types"
	"fmt"
	"log"
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
	status := kit.Request.FormValue("userStatus")

	log.Println("Payment methode", paymentMethod)
	specialRequest := kit.Request.FormValue("specialRequest")

	// Booking object
	booking := types.Bookings{
		UserID:         auth.UserID,
		CampID:         uint(campID),
		SpecialRequest: specialRequest,
		TotalPrice:     totalPrice,
		Status:         status,
		PaymentStatus:  "unpaid",
		PaymentMethod:  paymentMethod,
	}

	// Parse guests

	guests, err := GuestParsing(kit)

	if err != nil {
		return err
	}

	booking.Guests = guests

	//Parse services

	bookingServices, err := BookingServiceParsing(kit)

	if err != nil {
		return err
	}
	booking.Services = bookingServices

	// Create booking with nested associations
	if err := db.Get().Create(&booking).Error; err != nil {

		return err
	}

	return kit.Redirect(http.StatusSeeOther, "/")
}

func GuestParsing(kit *kit.Kit) ([]types.Guest, error) {
	var guests []types.Guest
	guestsCount, err := strconv.Atoi(kit.Request.FormValue("guestsCount"))
	if err != nil {
		return nil, err
	}
	for i := 0; i < guestsCount; i++ {
		guest := types.Guest{
			FirstName: kit.Request.FormValue(fmt.Sprintf("guests[%d][first_name]", i)),
			LastName:  kit.Request.FormValue(fmt.Sprintf("guests[%d][last_name]", i)),
			CIN:       kit.Request.FormValue(fmt.Sprintf("guests[%d][cin]", i)),
		}
		guests = append(guests, guest)
	}
	return guests, nil

}

func BookingServiceParsing(kit *kit.Kit) ([]types.BookingService, error) {
	var services []types.BookingService
	for k, v := range kit.Request.Form {
		if strings.HasPrefix(k, "service[") {
			idStr := strings.TrimSuffix(strings.TrimPrefix(k, "service["), "]")
			serviceID, _ := strconv.Atoi(idStr)
			qty, err := strconv.Atoi(v[0])
			if err != nil {
				return nil, err
			}
			if qty > 0 {
				services = append(services, types.BookingService{
					ServiceID: uint(serviceID),
					Quantity:  qty,
				})
			}
		}
	}

	return services, nil
}
