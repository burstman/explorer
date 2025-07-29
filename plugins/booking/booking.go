package booking

import (
	"explorer/app/db"
	"explorer/app/types"
	"explorer/plugins/auth"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/anthdm/superkit/kit"
)

func HandelCreateBooking(kit *kit.Kit) error {
	// Parse form
	if err := kit.Request.ParseForm(); err != nil {

		return err
	}

	// Parse scalar fields
	userID, err := getCurrentUserID(kit)
	if err != nil {
		return err
	}
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
		UserID:         uint(userID),
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

func getCurrentUserID(kit *kit.Kit) (uint, error) {
	sess := kit.GetSession("user-session")
	token, ok := sess.Values["sessionToken"]
	if !ok {
		log.Println("token not found")
		return 0, fmt.Errorf("token not found")
	}
	var session auth.Session
	err := db.Get().
		Preload("User").
		Find(&session, "token = ? AND expires_at > ?", token, time.Now()).Error
	if err != nil || session.ID == 0 {
		return 0, err
	}
	return session.UserID, nil
}
