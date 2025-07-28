package booking

import (
	"strconv"

	"github.com/anthdm/superkit/kit"
)

func HandelCreateBooking(kit *kit.Kit) error {
	err := kit.Request.ParseForm()
	if err != nil {
		return err
	}

	// Basic Fields
	campID, err := strconv.Atoi(r.FormValue("campID"))
	if err != nil {
		return err
	}
	bookingDate := kit.Request.FormValue("bookingDate")
	totalPrice, err := strconv.ParseFloat(kit.Request.FormValue("totalPrice"), 64)
	if err != nil {
		return err
	}

	paymentMethod := kit.Request.FormValue("payment_method")
	specialRequest := kit.Request.FormValue("specialRequest")

	return nil
}
